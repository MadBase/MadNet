// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "contracts/libraries/math/CryptoLibrary.sol";
import "contracts/interfaces/IETHDKG.sol";
import "contracts/interfaces/IETHDKGEvents.sol";
import "contracts/libraries/ethdkg/ETHDKGStorage.sol";
import "contracts/utils/ETHDKGUtils.sol";
import {ETHDKGErrorCodes} from "contracts/libraries/errorCodes/ETHDKGErrorCodes.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @custom:salt ETHDKGPhases
/// @custom:deploy-type deployUpgradeable
/// @custom:deploy-group ethdkg
/// @custom:deploy-group-index 1
contract ETHDKGPhases is ETHDKGStorage, IETHDKGEvents, ETHDKGUtils {
    using Strings for uint16;

    constructor() ETHDKGStorage() {}

    function register(uint256[2] memory publicKey) external {
        require(
            _ethdkgPhase == Phase.RegistrationOpen &&
                block.number >= _phaseStartBlock &&
                block.number < _phaseStartBlock + _phaseLength,
            ETHDKGErrorCodes.ETHDKG_NOT_IN_REGISTRATION_PHASE.toString()
        );
        require(
            publicKey[0] != 0 && publicKey[1] != 0,
            ETHDKGErrorCodes.ETHDKG_PUBLIC_KEY_ZERO.toString()
        );

        require(
            CryptoLibrary.bn128_is_on_curve(publicKey),
            ETHDKGErrorCodes.ETHDKG_PUBLIC_KEY_NOT_ON_CURVE.toString()
        );
        require(
            _participants[msg.sender].nonce < _nonce,
            ETHDKGErrorCodes.ETHDKG_PARTICIPANT_PARTICIPATING_IN_ROUND.toString()
        );
        uint32 numRegistered = uint32(_numParticipants);
        numRegistered++;
        _participants[msg.sender] = Participant({
            publicKey: publicKey,
            index: numRegistered,
            nonce: _nonce,
            phase: _ethdkgPhase,
            distributedSharesHash: 0x0,
            commitmentsFirstCoefficient: [uint256(0), uint256(0)],
            keyShares: [uint256(0), uint256(0)],
            gpkj: [uint256(0), uint256(0), uint256(0), uint256(0)]
        });

        emit AddressRegistered(msg.sender, numRegistered, _nonce, publicKey);
        if (
            _moveToNextPhase(
                Phase.ShareDistribution,
                IValidatorPool(_validatorPoolAddress()).getValidatorsCount(),
                numRegistered
            )
        ) {
            emit RegistrationComplete(block.number);
        }
    }

    function distributeShares(uint256[] memory encryptedShares, uint256[2][] memory commitments)
        external
    {
        require(
            _ethdkgPhase == Phase.ShareDistribution &&
                block.number >= _phaseStartBlock &&
                block.number < _phaseStartBlock + _phaseLength,
            ETHDKGErrorCodes.ETHDKG_NOT_IN_SHARED_DISTRIBUTION_PHASE.toString()
        );
        Participant memory participant = _participants[msg.sender];
        require(participant.nonce == _nonce, ETHDKGErrorCodes.ETHDKG_INVALID_NONCE.toString());
        require(
            participant.phase == Phase.RegistrationOpen,
            ETHDKGErrorCodes.ETHDKG_PARTICIPANT_DISTRIBUTED_SHARES_IN_ROUND.toString()
        );

        uint256 numValidators = IValidatorPool(_validatorPoolAddress()).getValidatorsCount();
        uint256 threshold = _getThreshold(numValidators);
        require(
            encryptedShares.length == numValidators - 1,
            ETHDKGErrorCodes.ETHDKG_INVALID_NUM_ENCRYPTED_SHARES.toString()
        );
        require(
            commitments.length == threshold + 1,
            ETHDKGErrorCodes.ETHDKG_INVALID_NUM_COMMITMENTS.toString()
        );
        for (uint256 k = 0; k <= threshold; k++) {
            require(
                CryptoLibrary.bn128_is_on_curve(commitments[k]),
                ETHDKGErrorCodes.ETHDKG_COMMITMENT_NOT_ON_CURVE.toString()
            );
            require(commitments[k][0] != 0, ETHDKGErrorCodes.ETHDKG_COMMITMENT_ZERO.toString());
        }

        bytes32 encryptedSharesHash = keccak256(abi.encodePacked(encryptedShares));
        bytes32 commitmentsHash = keccak256(abi.encodePacked(commitments));
        participant.distributedSharesHash = keccak256(
            abi.encodePacked(encryptedSharesHash, commitmentsHash)
        );
        require(
            participant.distributedSharesHash != 0x0,
            ETHDKGErrorCodes.ETHDKG_DISTRIBUTED_SHARE_HASH_ZERO.toString()
        );
        participant.commitmentsFirstCoefficient = commitments[0];
        participant.phase = Phase.ShareDistribution;

        _participants[msg.sender] = participant;
        uint256 numParticipants = _numParticipants + 1;

        emit SharesDistributed(
            msg.sender,
            participant.index,
            participant.nonce,
            encryptedShares,
            commitments
        );

        if (_moveToNextPhase(Phase.DisputeShareDistribution, numValidators, numParticipants)) {
            emit ShareDistributionComplete(block.number);
        }
    }

    function submitKeyShare(
        uint256[2] memory keyShareG1,
        uint256[2] memory keyShareG1CorrectnessProof,
        uint256[4] memory keyShareG2
    ) external {
        // Only progress if all participants distributed their shares
        // and no bad participant was found
        require(
            (_ethdkgPhase == Phase.KeyShareSubmission &&
                block.number >= _phaseStartBlock &&
                block.number < _phaseStartBlock + _phaseLength) ||
                (_ethdkgPhase == Phase.DisputeShareDistribution &&
                    block.number >= _phaseStartBlock + _phaseLength &&
                    block.number < _phaseStartBlock + 2 * _phaseLength &&
                    _badParticipants == 0),
            ETHDKGErrorCodes.ETHDKG_NOT_IN_KEYSHARE_SUBMISSION_PHASE.toString()
        );

        // Since we had a dispute stage prior this state we need to set global state in here
        if (_ethdkgPhase != Phase.KeyShareSubmission) {
            _setPhase(Phase.KeyShareSubmission);
        }
        Participant memory participant = _participants[msg.sender];
        require(
            participant.nonce == _nonce,
            ETHDKGErrorCodes.ETHDKG_KEYSHARE_PHASE_INVALID_NONCE.toString()
        );
        require(
            participant.phase == Phase.ShareDistribution,
            ETHDKGErrorCodes.ETHDKG_PARTICIPANT_SUBMITTED_KEYSHARES_IN_ROUND.toString()
        );

        require(
            CryptoLibrary.dleq_verify(
                [CryptoLibrary.H1x, CryptoLibrary.H1y],
                keyShareG1,
                [CryptoLibrary.G1x, CryptoLibrary.G1y],
                participant.commitmentsFirstCoefficient,
                keyShareG1CorrectnessProof
            ),
            ETHDKGErrorCodes.ETHDKG_INVALID_KEYSHARE_G1.toString()
        );
        require(
            CryptoLibrary.bn128_check_pairing(
                [
                    keyShareG1[0],
                    keyShareG1[1],
                    CryptoLibrary.H2xi,
                    CryptoLibrary.H2x,
                    CryptoLibrary.H2yi,
                    CryptoLibrary.H2y,
                    CryptoLibrary.H1x,
                    CryptoLibrary.H1y,
                    keyShareG2[0],
                    keyShareG2[1],
                    keyShareG2[2],
                    keyShareG2[3]
                ]
            ),
            ETHDKGErrorCodes.ETHDKG_INVALID_KEYSHARE_G2.toString()
        );

        participant.keyShares = keyShareG1;
        participant.phase = Phase.KeyShareSubmission;
        _participants[msg.sender] = participant;

        uint256[2] memory mpkG1 = _mpkG1;
        _mpkG1 = CryptoLibrary.bn128_add(
            [mpkG1[0], mpkG1[1], participant.keyShares[0], participant.keyShares[1]]
        );

        uint256 numParticipants = _numParticipants + 1;
        emit KeyShareSubmitted(
            msg.sender,
            participant.index,
            participant.nonce,
            keyShareG1,
            keyShareG1CorrectnessProof,
            keyShareG2
        );

        if (
            _moveToNextPhase(
                Phase.MPKSubmission,
                IValidatorPool(_validatorPoolAddress()).getValidatorsCount(),
                numParticipants
            )
        ) {
            emit KeyShareSubmissionComplete(block.number);
        }
    }

    function submitMasterPublicKey(uint256[4] memory masterPublicKey_) external {
        require(
            _ethdkgPhase == Phase.MPKSubmission &&
                block.number >= _phaseStartBlock &&
                block.number < _phaseStartBlock + _phaseLength,
            ETHDKGErrorCodes.ETHDKG_NOT_IN_MASTER_PUBLIC_KEY_SUBMISSION_PHASE.toString()
        );
        uint256[2] memory mpkG1 = _mpkG1;
        require(
            CryptoLibrary.bn128_check_pairing(
                [
                    mpkG1[0],
                    mpkG1[1],
                    CryptoLibrary.H2xi,
                    CryptoLibrary.H2x,
                    CryptoLibrary.H2yi,
                    CryptoLibrary.H2y,
                    CryptoLibrary.H1x,
                    CryptoLibrary.H1y,
                    masterPublicKey_[0],
                    masterPublicKey_[1],
                    masterPublicKey_[2],
                    masterPublicKey_[3]
                ]
            ),
            ETHDKGErrorCodes.ETHDKG_MASTER_PUBLIC_KEY_PAIRING_CHECK_FAILURE.toString()
        );

        _masterPublicKey = masterPublicKey_;

        _setPhase(Phase.GPKJSubmission);
        emit MPKSet(block.number, _nonce, masterPublicKey_);
    }

    function submitGPKJ(uint256[4] memory gpkj) external {
        //todo: should we evict all validators if no one sent the master public key in time?
        require(
            _ethdkgPhase == Phase.GPKJSubmission &&
                block.number >= _phaseStartBlock &&
                block.number < _phaseStartBlock + _phaseLength,
            ETHDKGErrorCodes.ETHDKG_NOT_IN_GPKJ_SUBMISSION_PHASE.toString()
        );

        Participant memory participant = _participants[msg.sender];

        require(
            participant.nonce == _nonce,
            ETHDKGErrorCodes.ETHDKG_KEYSHARE_PHASE_INVALID_NONCE.toString()
        );
        require(
            participant.phase == Phase.KeyShareSubmission,
            ETHDKGErrorCodes.ETHDKG_PARTICIPANT_SUBMITTED_GPKJ_IN_ROUND.toString()
        );

        require(
            gpkj[0] != 0 || gpkj[1] != 0 || gpkj[2] != 0 || gpkj[3] != 0,
            ETHDKGErrorCodes.ETHDKG_GPKJ_ZERO.toString()
        );

        participant.gpkj = gpkj;
        participant.phase = Phase.GPKJSubmission;
        _participants[msg.sender] = participant;

        emit ValidatorMemberAdded(
            msg.sender,
            participant.index,
            participant.nonce,
            ISnapshots(_snapshotsAddress()).getEpoch(),
            participant.gpkj[0],
            participant.gpkj[1],
            participant.gpkj[2],
            participant.gpkj[3]
        );

        uint256 numParticipants = _numParticipants + 1;
        if (
            _moveToNextPhase(
                Phase.DisputeGPKJSubmission,
                IValidatorPool(_validatorPoolAddress()).getValidatorsCount(),
                numParticipants
            )
        ) {
            emit GPKJSubmissionComplete(block.number);
        }
    }

    function complete() external {
        //todo: should we reward ppl here?
        require(
            (_ethdkgPhase == Phase.DisputeGPKJSubmission &&
                block.number >= _phaseStartBlock + _phaseLength) &&
                block.number < _phaseStartBlock + 2 * _phaseLength,
            ETHDKGErrorCodes.ETHDKG_NOT_IN_POST_GPKJ_DISPUTE_PHASE.toString()
        );
        require(_badParticipants == 0, ETHDKGErrorCodes.ETHDKG_REQUISITES_INCOMPLETE.toString());

        // Since we had a dispute stage prior this state we need to set global state in here
        _setPhase(Phase.Completion);

        IValidatorPool(_validatorPoolAddress()).completeETHDKG();

        uint256 epoch = ISnapshots(_snapshotsAddress()).getEpoch();
        uint256 ethHeight = ISnapshots(_snapshotsAddress()).getCommittedHeightFromLatestSnapshot();
        uint256 aliceNetHeight;
        if (_customAliceNetHeight == 0) {
            aliceNetHeight = ISnapshots(_snapshotsAddress()).getAliceNetHeightFromLatestSnapshot();
        } else {
            aliceNetHeight = _customAliceNetHeight;
            _customAliceNetHeight = 0;
        }
        emit ValidatorSetCompleted(
            uint8(IValidatorPool(_validatorPoolAddress()).getValidatorsCount()),
            _nonce,
            epoch,
            ethHeight,
            aliceNetHeight,
            _masterPublicKey[0],
            _masterPublicKey[1],
            _masterPublicKey[2],
            _masterPublicKey[3]
        );
    }

    function getMyAddress() public view returns (address) {
        return address(this);
    }

    function _setPhase(Phase phase_) internal {
        _ethdkgPhase = phase_;
        _phaseStartBlock = uint64(block.number);
        _numParticipants = 0;
    }

    function _moveToNextPhase(
        Phase phase_,
        uint256 numValidators_,
        uint256 numParticipants_
    ) internal returns (bool) {
        // if all validators have registered, we can proceed to the next phase
        if (numParticipants_ == numValidators_) {
            _setPhase(phase_);
            _phaseStartBlock += _confirmationLength;
            return true;
        } else {
            _numParticipants = uint32(numParticipants_);
            return false;
        }
    }

    function _isMasterPublicKeySet() internal view returns (bool) {
        return ((_masterPublicKey[0] != 0) ||
            (_masterPublicKey[1] != 0) ||
            (_masterPublicKey[2] != 0) ||
            (_masterPublicKey[3] != 0));
    }
}
