// Sources flattened with hardhat v2.9.1 https://hardhat.org

// File @openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (token/ERC20/IERC20.sol)

pragma solidity ^0.8.0;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IERC20Upgradeable {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender)
        external
        view
        returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `from` to `to` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );
}

// File @openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (token/ERC20/extensions/IERC20Metadata.sol)

pragma solidity ^0.8.0;

/**
 * @dev Interface for the optional metadata functions from the ERC20 standard.
 *
 * _Available since v4.1._
 */
interface IERC20MetadataUpgradeable is IERC20Upgradeable {
    /**
     * @dev Returns the name of the token.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the symbol of the token.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the decimals places of the token.
     */
    function decimals() external view returns (uint8);
}

// File @openzeppelin/contracts-upgradeable/utils/AddressUpgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (utils/Address.sol)

pragma solidity ^0.8.1;

/**
 * @dev Collection of functions related to the address type
 */
library AddressUpgradeable {
    /**
     * @dev Returns true if `account` is a contract.
     *
     * [IMPORTANT]
     * ====
     * It is unsafe to assume that an address for which this function returns
     * false is an externally-owned account (EOA) and not a contract.
     *
     * Among others, `isContract` will return false for the following
     * types of addresses:
     *
     *  - an externally-owned account
     *  - a contract in construction
     *  - an address where a contract will be created
     *  - an address where a contract lived, but was destroyed
     * ====
     *
     * [IMPORTANT]
     * ====
     * You shouldn't rely on `isContract` to protect against flash loan attacks!
     *
     * Preventing calls from contracts is highly discouraged. It breaks composability, breaks support for smart wallets
     * like Gnosis Safe, and does not provide security since it can be circumvented by calling from a contract
     * constructor.
     * ====
     */
    function isContract(address account) internal view returns (bool) {
        // This method relies on extcodesize/address.code.length, which returns 0
        // for contracts in construction, since the code is only stored at the end
        // of the constructor execution.

        return account.code.length > 0;
    }

    /**
     * @dev Replacement for Solidity's `transfer`: sends `amount` wei to
     * `recipient`, forwarding all available gas and reverting on errors.
     *
     * https://eips.ethereum.org/EIPS/eip-1884[EIP1884] increases the gas cost
     * of certain opcodes, possibly making contracts go over the 2300 gas limit
     * imposed by `transfer`, making them unable to receive funds via
     * `transfer`. {sendValue} removes this limitation.
     *
     * https://diligence.consensys.net/posts/2019/09/stop-using-soliditys-transfer-now/[Learn more].
     *
     * IMPORTANT: because control is transferred to `recipient`, care must be
     * taken to not create reentrancy vulnerabilities. Consider using
     * {ReentrancyGuard} or the
     * https://solidity.readthedocs.io/en/v0.5.11/security-considerations.html#use-the-checks-effects-interactions-pattern[checks-effects-interactions pattern].
     */
    function sendValue(address payable recipient, uint256 amount) internal {
        require(
            address(this).balance >= amount,
            "Address: insufficient balance"
        );

        (bool success, ) = recipient.call{value: amount}("");
        require(
            success,
            "Address: unable to send value, recipient may have reverted"
        );
    }

    /**
     * @dev Performs a Solidity function call using a low level `call`. A
     * plain `call` is an unsafe replacement for a function call: use this
     * function instead.
     *
     * If `target` reverts with a revert reason, it is bubbled up by this
     * function (like regular Solidity function calls).
     *
     * Returns the raw returned data. To convert to the expected return value,
     * use https://solidity.readthedocs.io/en/latest/units-and-global-variables.html?highlight=abi.decode#abi-encoding-and-decoding-functions[`abi.decode`].
     *
     * Requirements:
     *
     * - `target` must be a contract.
     * - calling `target` with `data` must not revert.
     *
     * _Available since v3.1._
     */
    function functionCall(address target, bytes memory data)
        internal
        returns (bytes memory)
    {
        return functionCall(target, data, "Address: low-level call failed");
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`], but with
     * `errorMessage` as a fallback revert reason when `target` reverts.
     *
     * _Available since v3.1._
     */
    function functionCall(
        address target,
        bytes memory data,
        string memory errorMessage
    ) internal returns (bytes memory) {
        return functionCallWithValue(target, data, 0, errorMessage);
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
     * but also transferring `value` wei to `target`.
     *
     * Requirements:
     *
     * - the calling contract must have an ETH balance of at least `value`.
     * - the called Solidity function must be `payable`.
     *
     * _Available since v3.1._
     */
    function functionCallWithValue(
        address target,
        bytes memory data,
        uint256 value
    ) internal returns (bytes memory) {
        return
            functionCallWithValue(
                target,
                data,
                value,
                "Address: low-level call with value failed"
            );
    }

    /**
     * @dev Same as {xref-Address-functionCallWithValue-address-bytes-uint256-}[`functionCallWithValue`], but
     * with `errorMessage` as a fallback revert reason when `target` reverts.
     *
     * _Available since v3.1._
     */
    function functionCallWithValue(
        address target,
        bytes memory data,
        uint256 value,
        string memory errorMessage
    ) internal returns (bytes memory) {
        require(
            address(this).balance >= value,
            "Address: insufficient balance for call"
        );
        require(isContract(target), "Address: call to non-contract");

        (bool success, bytes memory returndata) = target.call{value: value}(
            data
        );
        return verifyCallResult(success, returndata, errorMessage);
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
     * but performing a static call.
     *
     * _Available since v3.3._
     */
    function functionStaticCall(address target, bytes memory data)
        internal
        view
        returns (bytes memory)
    {
        return
            functionStaticCall(
                target,
                data,
                "Address: low-level static call failed"
            );
    }

    /**
     * @dev Same as {xref-Address-functionCall-address-bytes-string-}[`functionCall`],
     * but performing a static call.
     *
     * _Available since v3.3._
     */
    function functionStaticCall(
        address target,
        bytes memory data,
        string memory errorMessage
    ) internal view returns (bytes memory) {
        require(isContract(target), "Address: static call to non-contract");

        (bool success, bytes memory returndata) = target.staticcall(data);
        return verifyCallResult(success, returndata, errorMessage);
    }

    /**
     * @dev Tool to verifies that a low level call was successful, and revert if it wasn't, either by bubbling the
     * revert reason using the provided one.
     *
     * _Available since v4.3._
     */
    function verifyCallResult(
        bool success,
        bytes memory returndata,
        string memory errorMessage
    ) internal pure returns (bytes memory) {
        if (success) {
            return returndata;
        } else {
            // Look for revert reason and bubble it up if present
            if (returndata.length > 0) {
                // The easiest way to bubble the revert reason is using memory via assembly

                assembly {
                    let returndata_size := mload(returndata)
                    revert(add(32, returndata), returndata_size)
                }
            } else {
                revert(errorMessage);
            }
        }
    }
}

// File @openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (proxy/utils/Initializable.sol)

pragma solidity ^0.8.0;

/**
 * @dev This is a base contract to aid in writing upgradeable contracts, or any kind of contract that will be deployed
 * behind a proxy. Since proxied contracts do not make use of a constructor, it's common to move constructor logic to an
 * external initializer function, usually called `initialize`. It then becomes necessary to protect this initializer
 * function so it can only be called once. The {initializer} modifier provided by this contract will have this effect.
 *
 * TIP: To avoid leaving the proxy in an uninitialized state, the initializer function should be called as early as
 * possible by providing the encoded function call as the `_data` argument to {ERC1967Proxy-constructor}.
 *
 * CAUTION: When used with inheritance, manual care must be taken to not invoke a parent initializer twice, or to ensure
 * that all initializers are idempotent. This is not verified automatically as constructors are by Solidity.
 *
 * [CAUTION]
 * ====
 * Avoid leaving a contract uninitialized.
 *
 * An uninitialized contract can be taken over by an attacker. This applies to both a proxy and its implementation
 * contract, which may impact the proxy. To initialize the implementation contract, you can either invoke the
 * initializer manually, or you can include a constructor to automatically mark it as initialized when it is deployed:
 *
 * [.hljs-theme-light.nopadding]
 * ```
 * /// @custom:oz-upgrades-unsafe-allow constructor
 * constructor() initializer {}
 * ```
 * ====
 */
abstract contract Initializable {
    /**
     * @dev Indicates that the contract has been initialized.
     */
    bool private _initialized;

    /**
     * @dev Indicates that the contract is in the process of being initialized.
     */
    bool private _initializing;

    /**
     * @dev Modifier to protect an initializer function from being invoked twice.
     */
    modifier initializer() {
        // If the contract is initializing we ignore whether _initialized is set in order to support multiple
        // inheritance patterns, but we only do this in the context of a constructor, because in other contexts the
        // contract may have been reentered.
        require(
            _initializing ? _isConstructor() : !_initialized,
            "Initializable: contract is already initialized"
        );

        bool isTopLevelCall = !_initializing;
        if (isTopLevelCall) {
            _initializing = true;
            _initialized = true;
        }

        _;

        if (isTopLevelCall) {
            _initializing = false;
        }
    }

    /**
     * @dev Modifier to protect an initialization function so that it can only be invoked by functions with the
     * {initializer} modifier, directly or indirectly.
     */
    modifier onlyInitializing() {
        require(_initializing, "Initializable: contract is not initializing");
        _;
    }

    function _isConstructor() private view returns (bool) {
        return !AddressUpgradeable.isContract(address(this));
    }
}

// File @openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/Context.sol)

pragma solidity ^0.8.0;

/**
 * @dev Provides information about the current execution context, including the
 * sender of the transaction and its data. While these are generally available
 * via msg.sender and msg.data, they should not be accessed in such a direct
 * manner, since when dealing with meta-transactions the account sending and
 * paying for execution may not be the actual sender (as far as an application
 * is concerned).
 *
 * This contract is only required for intermediate, library-like contracts.
 */
abstract contract ContextUpgradeable is Initializable {
    function __Context_init() internal onlyInitializing {}

    function __Context_init_unchained() internal onlyInitializing {}

    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }

    /**
     * @dev This empty reserved space is put in place to allow future versions to add new
     * variables without shifting down storage in the inheritance chain.
     * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
     */
    uint256[50] private __gap;
}

// File @openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol@v4.5.2

// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.5.0) (token/ERC20/ERC20.sol)

pragma solidity ^0.8.0;

/**
 * @dev Implementation of the {IERC20} interface.
 *
 * This implementation is agnostic to the way tokens are created. This means
 * that a supply mechanism has to be added in a derived contract using {_mint}.
 * For a generic mechanism see {ERC20PresetMinterPauser}.
 *
 * TIP: For a detailed writeup see our guide
 * https://forum.zeppelin.solutions/t/how-to-implement-erc20-supply-mechanisms/226[How
 * to implement supply mechanisms].
 *
 * We have followed general OpenZeppelin Contracts guidelines: functions revert
 * instead returning `false` on failure. This behavior is nonetheless
 * conventional and does not conflict with the expectations of ERC20
 * applications.
 *
 * Additionally, an {Approval} event is emitted on calls to {transferFrom}.
 * This allows applications to reconstruct the allowance for all accounts just
 * by listening to said events. Other implementations of the EIP may not emit
 * these events, as it isn't required by the specification.
 *
 * Finally, the non-standard {decreaseAllowance} and {increaseAllowance}
 * functions have been added to mitigate the well-known issues around setting
 * allowances. See {IERC20-approve}.
 */
contract ERC20Upgradeable is
    Initializable,
    ContextUpgradeable,
    IERC20Upgradeable,
    IERC20MetadataUpgradeable
{
    mapping(address => uint256) private _balances;

    mapping(address => mapping(address => uint256)) private _allowances;

    uint256 private _totalSupply;

    string private _name;
    string private _symbol;

    /**
     * @dev Sets the values for {name} and {symbol}.
     *
     * The default value of {decimals} is 18. To select a different value for
     * {decimals} you should overload it.
     *
     * All two of these values are immutable: they can only be set once during
     * construction.
     */
    function __ERC20_init(string memory name_, string memory symbol_)
        internal
        onlyInitializing
    {
        __ERC20_init_unchained(name_, symbol_);
    }

    function __ERC20_init_unchained(string memory name_, string memory symbol_)
        internal
        onlyInitializing
    {
        _name = name_;
        _symbol = symbol_;
    }

    /**
     * @dev Returns the name of the token.
     */
    function name() public view virtual override returns (string memory) {
        return _name;
    }

    /**
     * @dev Returns the symbol of the token, usually a shorter version of the
     * name.
     */
    function symbol() public view virtual override returns (string memory) {
        return _symbol;
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5.05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei. This is the value {ERC20} uses, unless this function is
     * overridden;
     *
     * NOTE: This information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * {IERC20-balanceOf} and {IERC20-transfer}.
     */
    function decimals() public view virtual override returns (uint8) {
        return 18;
    }

    /**
     * @dev See {IERC20-totalSupply}.
     */
    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev See {IERC20-balanceOf}.
     */
    function balanceOf(address account)
        public
        view
        virtual
        override
        returns (uint256)
    {
        return _balances[account];
    }

    /**
     * @dev See {IERC20-transfer}.
     *
     * Requirements:
     *
     * - `to` cannot be the zero address.
     * - the caller must have a balance of at least `amount`.
     */
    function transfer(address to, uint256 amount)
        public
        virtual
        override
        returns (bool)
    {
        address owner = _msgSender();
        _transfer(owner, to, amount);
        return true;
    }

    /**
     * @dev See {IERC20-allowance}.
     */
    function allowance(address owner, address spender)
        public
        view
        virtual
        override
        returns (uint256)
    {
        return _allowances[owner][spender];
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * NOTE: If `amount` is the maximum `uint256`, the allowance is not updated on
     * `transferFrom`. This is semantically equivalent to an infinite approval.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount)
        public
        virtual
        override
        returns (bool)
    {
        address owner = _msgSender();
        _approve(owner, spender, amount);
        return true;
    }

    /**
     * @dev See {IERC20-transferFrom}.
     *
     * Emits an {Approval} event indicating the updated allowance. This is not
     * required by the EIP. See the note at the beginning of {ERC20}.
     *
     * NOTE: Does not update the allowance if the current allowance
     * is the maximum `uint256`.
     *
     * Requirements:
     *
     * - `from` and `to` cannot be the zero address.
     * - `from` must have a balance of at least `amount`.
     * - the caller must have allowance for ``from``'s tokens of at least
     * `amount`.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) public virtual override returns (bool) {
        address spender = _msgSender();
        _spendAllowance(from, spender, amount);
        _transfer(from, to, amount);
        return true;
    }

    /**
     * @dev Atomically increases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function increaseAllowance(address spender, uint256 addedValue)
        public
        virtual
        returns (bool)
    {
        address owner = _msgSender();
        _approve(owner, spender, _allowances[owner][spender] + addedValue);
        return true;
    }

    /**
     * @dev Atomically decreases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     * - `spender` must have allowance for the caller of at least
     * `subtractedValue`.
     */
    function decreaseAllowance(address spender, uint256 subtractedValue)
        public
        virtual
        returns (bool)
    {
        address owner = _msgSender();
        uint256 currentAllowance = _allowances[owner][spender];
        require(
            currentAllowance >= subtractedValue,
            "ERC20: decreased allowance below zero"
        );
        unchecked {
            _approve(owner, spender, currentAllowance - subtractedValue);
        }

        return true;
    }

    /**
     * @dev Moves `amount` of tokens from `sender` to `recipient`.
     *
     * This internal function is equivalent to {transfer}, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a {Transfer} event.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `from` must have a balance of at least `amount`.
     */
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");

        _beforeTokenTransfer(from, to, amount);

        uint256 fromBalance = _balances[from];
        require(
            fromBalance >= amount,
            "ERC20: transfer amount exceeds balance"
        );
        unchecked {
            _balances[from] = fromBalance - amount;
        }
        _balances[to] += amount;

        emit Transfer(from, to, amount);

        _afterTokenTransfer(from, to, amount);
    }

    /** @dev Creates `amount` tokens and assigns them to `account`, increasing
     * the total supply.
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     */
    function _mint(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: mint to the zero address");

        _beforeTokenTransfer(address(0), account, amount);

        _totalSupply += amount;
        _balances[account] += amount;
        emit Transfer(address(0), account, amount);

        _afterTokenTransfer(address(0), account, amount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`, reducing the
     * total supply.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     * - `account` must have at least `amount` tokens.
     */
    function _burn(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: burn from the zero address");

        _beforeTokenTransfer(account, address(0), amount);

        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "ERC20: burn amount exceeds balance");
        unchecked {
            _balances[account] = accountBalance - amount;
        }
        _totalSupply -= amount;

        emit Transfer(account, address(0), amount);

        _afterTokenTransfer(account, address(0), amount);
    }

    /**
     * @dev Sets `amount` as the allowance of `spender` over the `owner` s tokens.
     *
     * This internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     * Requirements:
     *
     * - `owner` cannot be the zero address.
     * - `spender` cannot be the zero address.
     */
    function _approve(
        address owner,
        address spender,
        uint256 amount
    ) internal virtual {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    /**
     * @dev Spend `amount` form the allowance of `owner` toward `spender`.
     *
     * Does not update the allowance amount in case of infinite allowance.
     * Revert if not enough allowance is available.
     *
     * Might emit an {Approval} event.
     */
    function _spendAllowance(
        address owner,
        address spender,
        uint256 amount
    ) internal virtual {
        uint256 currentAllowance = allowance(owner, spender);
        if (currentAllowance != type(uint256).max) {
            require(
                currentAllowance >= amount,
                "ERC20: insufficient allowance"
            );
            unchecked {
                _approve(owner, spender, currentAllowance - amount);
            }
        }
    }

    /**
     * @dev Hook that is called before any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * will be transferred to `to`.
     * - when `from` is zero, `amount` tokens will be minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens will be burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {}

    /**
     * @dev Hook that is called after any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * has been transferred to `to`.
     * - when `from` is zero, `amount` tokens have been minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens have been burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _afterTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {}

    /**
     * @dev This empty reserved space is put in place to allow future versions to add new
     * variables without shifting down storage in the inheritance chain.
     * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
     */
    uint256[45] private __gap;
}

// File contracts/libraries/errorCodes/AdminErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library AdminErrorCodes {
    // AdminErrorCodes error codes
    bytes32 public constant ADMIN_SENDER_MUST_BE_ADMIN = "1700"; //"Must be admin"
}

// File contracts/utils/Admin.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract Admin {
    // _admin is a privileged role
    address internal _admin;

    /// @dev onlyAdmin enforces msg.sender is _admin
    modifier onlyAdmin() {
        require(
            msg.sender == _admin,
            string(abi.encodePacked(AdminErrorCodes.ADMIN_SENDER_MUST_BE_ADMIN))
        );
        _;
    }

    constructor(address admin_) {
        _admin = admin_;
    }

    /// @dev assigns a new admin may only be called by _admin
    function setAdmin(address admin_) public virtual onlyAdmin {
        _setAdmin(admin_);
    }

    /// @dev getAdmin returns the current _admin
    function getAdmin() public view returns (address) {
        return _admin;
    }

    // assigns a new admin may only be called by _admin
    function _setAdmin(address admin_) internal {
        _admin = admin_;
    }
}

// File contracts/libraries/errorCodes/MutexErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library MutexErrorCodes {
    // Mutex error codes
    bytes32 public constant MUTEX_LOCKED = "2300"; //"Mutex: Couldn't acquire the lock!"
}

// File contracts/utils/Mutex.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract Mutex {
    uint256 internal constant _LOCKED = 1;
    uint256 internal constant _UNLOCKED = 2;
    uint256 internal _mutex;

    modifier withLock() {
        require(
            _mutex != _LOCKED,
            string(abi.encodePacked(MutexErrorCodes.MUTEX_LOCKED))
        );
        _mutex = _LOCKED;
        _;
        _mutex = _UNLOCKED;
    }

    constructor() {
        _mutex = _UNLOCKED;
    }
}

// File contracts/libraries/errorCodes/MagicValueErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library MagicValueErrorCodes {
    // MagicValue error codes
    bytes32 public constant MAGICVALUE_BAD_MAGIC = "2200"; //"BAD MAGIC"
}

// File contracts/utils/MagicValue.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract MagicValue {
    // _MAGIC_VALUE is a constant that may be used to prevent
    // a user from calling a dangerous method without significant
    // effort or ( hopefully ) reading the code to understand the risk
    uint8 internal constant _MAGIC_VALUE = 42;

    modifier checkMagic(uint8 magic_) {
        require(
            magic_ == _getMagic(),
            string(abi.encodePacked(MagicValueErrorCodes.MAGICVALUE_BAD_MAGIC))
        );
        _;
    }

    // _getMagic returns the magic constant
    function _getMagic() internal pure returns (uint8) {
        return _MAGIC_VALUE;
    }
}

// File contracts/interfaces/IMagicEthTransfer.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

interface IMagicEthTransfer {
    function depositEth(uint8 magic_) external payable;
}

// File contracts/utils/MagicEthTransfer.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract MagicEthTransfer is MagicValue {
    function _safeTransferEthWithMagic(IMagicEthTransfer to_, uint256 amount_)
        internal
    {
        to_.depositEth{value: amount_}(_getMagic());
    }
}

// File contracts/utils/EthSafeTransfer.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract EthSafeTransfer {
    /// @notice _safeTransferEth performs a transfer of Eth using the call
    /// method / this function is resistant to breaking gas price changes and /
    /// performs call in a safe manner by reverting on failure. / this function
    /// will return without performing a call or reverting, / if amount_ is zero
    function _safeTransferEth(address to_, uint256 amount_) internal {
        if (amount_ == 0) {
            return;
        }
        require(
            to_ != address(0),
            "EthSafeTransfer: cannot transfer ETH to address 0x0"
        );
        address payable caller = payable(to_);
        (bool success, ) = caller.call{value: amount_}("");
        require(success, "EthSafeTransfer: Transfer failed.");
    }
}

// File contracts/libraries/math/Sigmoid.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.0;

abstract contract Sigmoid {
    function _fx(uint256 x) internal pure returns (uint256) {
        return
            201 *
            x +
            504975246918103897890400 -
            _sqrt(
                200**2 *
                    ((_safeAbsSub(2500000000000000000000, x))**2 +
                        125000000000000005611050234958650739260304)
            );
    }

    function _fp(uint256 p) internal pure returns (uint256) {
        return
            (201 *
                p +
                _sqrt(
                    200**2 *
                        (p**2 +
                            1015056498152694417606544040542564448516855541904 -
                            2014950493836207795780800 *
                            p)
                ) -
                201500024630538883475970400) / 401;
    }

    function _safeAbsSub(uint256 a, uint256 b) internal pure returns (uint256) {
        return _max(a, b) - _min(a, b);
    }

    function _min(uint256 a_, uint256 b_) internal pure returns (uint256) {
        if (a_ <= b_) {
            return a_;
        }
        return b_;
    }

    function _max(uint256 a_, uint256 b_) internal pure returns (uint256) {
        if (a_ >= b_) {
            return a_;
        }
        return b_;
    }

    /// @notice Calculates the square root of x, rounding down.
    /// @dev Uses the Babylonian method https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Babylonian_method.
    ///
    /// Caveats:
    /// - This function does not work with fixed-point numbers.
    ///
    /// @param x The uint256 number for which to calculate the square root.
    /// @return result The result as an uint256.
    function _sqrt(uint256 x) internal pure returns (uint256 result) {
        if (x == 0) {
            return 0;
        }

        // Set the initial guess to the closest power of two that is higher than x.
        uint256 xAux = uint256(x);
        result = 1;
        if (xAux >= 0x100000000000000000000000000000000) {
            xAux >>= 128;
            result <<= 64;
        }
        if (xAux >= 0x10000000000000000) {
            xAux >>= 64;
            result <<= 32;
        }
        if (xAux >= 0x100000000) {
            xAux >>= 32;
            result <<= 16;
        }
        if (xAux >= 0x10000) {
            xAux >>= 16;
            result <<= 8;
        }
        if (xAux >= 0x100) {
            xAux >>= 8;
            result <<= 4;
        }
        if (xAux >= 0x10) {
            xAux >>= 4;
            result <<= 2;
        }
        if (xAux >= 0x8) {
            result <<= 1;
        }

        // The operations can never overflow because the result is max 2^127 when it enters this block.
        unchecked {
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1;
            result = (result + x / result) >> 1; // Seven iterations should be enough
            uint256 roundedDownResult = x / result;
            return result >= roundedDownResult ? roundedDownResult : result;
        }
    }
}

// File contracts/utils/DeterministicAddress.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

abstract contract DeterministicAddress {
    function getMetamorphicContractAddress(bytes32 _salt, address _factory)
        public
        pure
        returns (address)
    {
        // byte code for metamorphic contract
        // 6020363636335afa1536363636515af43d36363e3d36f3
        bytes32 metamorphicContractBytecodeHash_ = 0x1c0bf703a3415cada9785e89e9d70314c3111ae7d8e04f33bb42eb1d264088be;
        return
            address(
                uint160(
                    uint256(
                        keccak256(
                            abi.encodePacked(
                                hex"ff",
                                _factory,
                                _salt,
                                metamorphicContractBytecodeHash_
                            )
                        )
                    )
                )
            );
    }
}

// File contracts/libraries/errorCodes/ImmutableAuthErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library ImmutableAuthErrorCodes {
    // ImmutableAuth error codes
    bytes32 public constant IMMUTEABLEAUTH_ONLY_FACTORY = "2000"; //"onlyFactory"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKEN = "2001"; //"onlyAToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_FOUNDATION = "2002"; //"onlyFoundation"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_GOVERNANCE = "2003"; // "onlyGovernance"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_LIQUIDITYPROVIDERSTAKING =
        "2004"; // "onlyLiquidityProviderStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_BTOKEN = "2005"; // "onlyBToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_MADTOKEN = "2006"; // "onlyMadToken"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_PUBLICSTAKING = "2007"; // "onlyPublicStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_SNAPSHOTS = "2008"; // "onlySnapshots"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_STAKINGPOSITIONDESCRIPTOR =
        "2009"; // "onlyStakingPositionDescriptor"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_VALIDATORPOOL = "2010"; // "onlyValidatorPool"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_VALIDATORSTAKING = "2011"; // "onlyValidatorStaking"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKENBURNER = "2012"; // "onlyATokenBurner"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ATOKENMINTER = "2013"; // "onlyATokenMinter"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKGACCUSATIONS = "2014"; // "onlyETHDKGAccusations"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKGPHASES = "2015"; // "onlyETHDKGPhases"
    bytes32 public constant IMMUTEABLEAUTH_ONLY_ETHDKG = "2016"; // "onlyETHDKG"
}

// File contracts/utils/ImmutableAuth.sol

// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;

    modifier onlyFactory() {
        require(
            msg.sender == _factory,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FACTORY
                )
            )
        );
        _;
    }

    constructor(address factory_) {
        _factory = factory_;
    }

    function _factoryAddress() internal view returns (address) {
        return _factory;
    }
}

abstract contract ImmutableAToken is ImmutableFactory {
    address private immutable _aToken;

    modifier onlyAToken() {
        require(
            msg.sender == _aToken,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKEN
                )
            )
        );
        _;
    }

    constructor() {
        _aToken = getMetamorphicContractAddress(
            0x41546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenAddress() internal view returns (address) {
        return _aToken;
    }

    function _saltForAToken() internal pure returns (bytes32) {
        return
            0x41546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenBurner is ImmutableFactory {
    address private immutable _aTokenBurner;

    modifier onlyATokenBurner() {
        require(
            msg.sender == _aTokenBurner,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENBURNER
                )
            )
        );
        _;
    }

    constructor() {
        _aTokenBurner = getMetamorphicContractAddress(
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenBurnerAddress() internal view returns (address) {
        return _aTokenBurner;
    }

    function _saltForATokenBurner() internal pure returns (bytes32) {
        return
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenMinter is ImmutableFactory {
    address private immutable _aTokenMinter;

    modifier onlyATokenMinter() {
        require(
            msg.sender == _aTokenMinter,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENMINTER
                )
            )
        );
        _;
    }

    constructor() {
        _aTokenMinter = getMetamorphicContractAddress(
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenMinterAddress() internal view returns (address) {
        return _aTokenMinter;
    }

    function _saltForATokenMinter() internal pure returns (bytes32) {
        return
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableBToken is ImmutableFactory {
    address private immutable _bToken;

    modifier onlyBToken() {
        require(
            msg.sender == _bToken,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_BTOKEN
                )
            )
        );
        _;
    }

    constructor() {
        _bToken = getMetamorphicContractAddress(
            0x42546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _bTokenAddress() internal view returns (address) {
        return _bToken;
    }

    function _saltForBToken() internal pure returns (bytes32) {
        return
            0x42546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableFoundation is ImmutableFactory {
    address private immutable _foundation;

    modifier onlyFoundation() {
        require(
            msg.sender == _foundation,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FOUNDATION
                )
            )
        );
        _;
    }

    constructor() {
        _foundation = getMetamorphicContractAddress(
            0x466f756e646174696f6e00000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _foundationAddress() internal view returns (address) {
        return _foundation;
    }

    function _saltForFoundation() internal pure returns (bytes32) {
        return
            0x466f756e646174696f6e00000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableGovernance is ImmutableFactory {
    address private immutable _governance;

    modifier onlyGovernance() {
        require(
            msg.sender == _governance,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_GOVERNANCE
                )
            )
        );
        _;
    }

    constructor() {
        _governance = getMetamorphicContractAddress(
            0x476f7665726e616e636500000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _governanceAddress() internal view returns (address) {
        return _governance;
    }

    function _saltForGovernance() internal pure returns (bytes32) {
        return
            0x476f7665726e616e636500000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableLiquidityProviderStaking is ImmutableFactory {
    address private immutable _liquidityProviderStaking;

    modifier onlyLiquidityProviderStaking() {
        require(
            msg.sender == _liquidityProviderStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_LIQUIDITYPROVIDERSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _liquidityProviderStaking = getMetamorphicContractAddress(
            0x4c697175696469747950726f76696465725374616b696e670000000000000000,
            _factoryAddress()
        );
    }

    function _liquidityProviderStakingAddress()
        internal
        view
        returns (address)
    {
        return _liquidityProviderStaking;
    }

    function _saltForLiquidityProviderStaking()
        internal
        pure
        returns (bytes32)
    {
        return
            0x4c697175696469747950726f76696465725374616b696e670000000000000000;
    }
}

abstract contract ImmutablePublicStaking is ImmutableFactory {
    address private immutable _publicStaking;

    modifier onlyPublicStaking() {
        require(
            msg.sender == _publicStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_PUBLICSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _publicStaking = getMetamorphicContractAddress(
            0x5075626c69635374616b696e6700000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _publicStakingAddress() internal view returns (address) {
        return _publicStaking;
    }

    function _saltForPublicStaking() internal pure returns (bytes32) {
        return
            0x5075626c69635374616b696e6700000000000000000000000000000000000000;
    }
}

abstract contract ImmutableSnapshots is ImmutableFactory {
    address private immutable _snapshots;

    modifier onlySnapshots() {
        require(
            msg.sender == _snapshots,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_SNAPSHOTS
                )
            )
        );
        _;
    }

    constructor() {
        _snapshots = getMetamorphicContractAddress(
            0x536e617073686f74730000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _snapshotsAddress() internal view returns (address) {
        return _snapshots;
    }

    function _saltForSnapshots() internal pure returns (bytes32) {
        return
            0x536e617073686f74730000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakingPositionDescriptor is ImmutableFactory {
    address private immutable _stakingPositionDescriptor;

    modifier onlyStakingPositionDescriptor() {
        require(
            msg.sender == _stakingPositionDescriptor,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_STAKINGPOSITIONDESCRIPTOR
                )
            )
        );
        _;
    }

    constructor() {
        _stakingPositionDescriptor = getMetamorphicContractAddress(
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000,
            _factoryAddress()
        );
    }

    function _stakingPositionDescriptorAddress()
        internal
        view
        returns (address)
    {
        return _stakingPositionDescriptor;
    }

    function _saltForStakingPositionDescriptor()
        internal
        pure
        returns (bytes32)
    {
        return
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000;
    }
}

abstract contract ImmutableValidatorPool is ImmutableFactory {
    address private immutable _validatorPool;

    modifier onlyValidatorPool() {
        require(
            msg.sender == _validatorPool,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORPOOL
                )
            )
        );
        _;
    }

    constructor() {
        _validatorPool = getMetamorphicContractAddress(
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorPoolAddress() internal view returns (address) {
        return _validatorPool;
    }

    function _saltForValidatorPool() internal pure returns (bytes32) {
        return
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000;
    }
}

abstract contract ImmutableValidatorStaking is ImmutableFactory {
    address private immutable _validatorStaking;

    modifier onlyValidatorStaking() {
        require(
            msg.sender == _validatorStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _validatorStaking = getMetamorphicContractAddress(
            0x56616c696461746f725374616b696e6700000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorStakingAddress() internal view returns (address) {
        return _validatorStaking;
    }

    function _saltForValidatorStaking() internal pure returns (bytes32) {
        return
            0x56616c696461746f725374616b696e6700000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGAccusations is ImmutableFactory {
    address private immutable _ethdkgAccusations;

    modifier onlyETHDKGAccusations() {
        require(
            msg.sender == _ethdkgAccusations,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes
                        .IMMUTEABLEAUTH_ONLY_ETHDKGACCUSATIONS
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkgAccusations = getMetamorphicContractAddress(
            0x455448444b4741636375736174696f6e73000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAccusationsAddress() internal view returns (address) {
        return _ethdkgAccusations;
    }

    function _saltForETHDKGAccusations() internal pure returns (bytes32) {
        return
            0x455448444b4741636375736174696f6e73000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGPhases is ImmutableFactory {
    address private immutable _ethdkgPhases;

    modifier onlyETHDKGPhases() {
        require(
            msg.sender == _ethdkgPhases,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKGPHASES
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkgPhases = getMetamorphicContractAddress(
            0x455448444b475068617365730000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgPhasesAddress() internal view returns (address) {
        return _ethdkgPhases;
    }

    function _saltForETHDKGPhases() internal pure returns (bytes32) {
        return
            0x455448444b475068617365730000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKG is ImmutableFactory {
    address private immutable _ethdkg;

    modifier onlyETHDKG() {
        require(
            msg.sender == _ethdkg,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKG
                )
            )
        );
        _;
    }

    constructor() {
        _ethdkg = getMetamorphicContractAddress(
            0x455448444b470000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAddress() internal view returns (address) {
        return _ethdkg;
    }

    function _saltForETHDKG() internal pure returns (bytes32) {
        return
            0x455448444b470000000000000000000000000000000000000000000000000000;
    }
}

// File contracts/libraries/errorCodes/BTokenErrorCodes.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

library BTokenErrorCodes {
    // BToken error codes
    bytes32 public constant BTOKEN_INVALID_DEPOSIT_ID = "300"; //"BToken: Invalid deposit ID!"
    bytes32 public constant BTOKEN_INVALID_BALANCE = "301"; //"BToken: Address balance should be always greater than the pool balance!"
    bytes32 public constant BTOKEN_INVALID_BURN_AMOUNT = "302"; //"BToken: The number of BTokens to be burn should be greater than 0!"
    bytes32 public constant BTOKEN_CONTRACTS_DISALLOWED_DEPOSITS = "303"; //"BToken: Contracts cannot make BTokens deposits!"
    bytes32 public constant BTOKEN_DEPOSIT_AMOUNT_ZERO = "304"; //"BToken: The deposit amount must be greater than zero!"
    bytes32 public constant BTOKEN_DEPOSIT_BURN_FAIL = "305"; //"BToken: Burn failed during the deposit!"
    bytes32 public constant BTOKEN_MARKET_SPREAD_TOO_LOW = "306"; //"BToken: requires at least 4 WEI"
    bytes32 public constant BTOKEN_MINT_INSUFFICIENT_ETH = "307"; //"BToken: could not mint deposit with minimum BTokens given the ether sent!"
    bytes32 public constant BTOKEN_MINIMUM_MINT_NOT_MET = "308"; //"BToken: could not mint minimum BTokens"
    bytes32 public constant BTOKEN_MINIMUM_BURN_NOT_MET = "309"; //"BToken: Couldn't burn the minEth amount"
    bytes32 public constant BTOKEN_SPLIT_VALUE_SUM_ERROR = "310"; //"BToken: All the split values must sum to _PERCENTAGE_SCALE!"
    bytes32 public constant BTOKEN_BURN_AMOUNT_EXCEEDS_SUPPLY = "311"; //"BToken: The number of tokens to be burned is greater than the Total Supply!"
}

// File contracts/BToken.sol

// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

/// @custom:salt BToken
/// @custom:deploy-type deployStatic
contract BToken is
    ERC20Upgradeable,
    Admin,
    Mutex,
    MagicEthTransfer,
    EthSafeTransfer,
    Sigmoid,
    ImmutableFactory,
    ImmutablePublicStaking,
    ImmutableValidatorStaking,
    ImmutableLiquidityProviderStaking,
    ImmutableFoundation
{
    struct Deposit {
        uint8 accountType;
        address account;
        uint256 value;
    }

    // multiply factor for the selling/minting bonding curve
    uint256 internal constant _MARKET_SPREAD = 4;

    // Scaling factor to get the staking percentages
    uint256 internal constant _PERCENTAGE_SCALE = 1000;

    // Balance in ether that is hold in the contract after minting and burning
    uint256 internal _poolBalance;

    // Value of the percentages that will send to each staking contract. Divide
    // this value by _PERCENTAGE_SCALE = 1000 to get the corresponding percentages.
    // These values must sum to 1000.
    uint256 internal _validatorStakingSplit;
    uint256 internal _publicStakingSplit;
    uint256 internal _liquidityProviderStakingSplit;
    uint256 internal _protocolFee;

    // Monotonically increasing variable to track the BTokens deposits.
    uint256 internal _depositID;

    // Total amount of BTokens that were deposited in the AliceNet chain. The
    // BTokens deposited in the AliceNet are burned by this contract.
    uint256 internal _totalDeposited;

    // Tracks the amount of each deposit. Key is deposit id, value is amount
    // deposited.
    mapping(uint256 => Deposit) internal _deposits;

    /// @notice Event emitted when a deposit is received
    event DepositReceived(
        uint256 indexed depositID,
        uint8 indexed accountType,
        address indexed depositor,
        uint256 amount
    );

    constructor()
        Admin(msg.sender)
        Mutex()
        ImmutableFactory(msg.sender)
        ImmutablePublicStaking()
        ImmutableValidatorStaking()
        ImmutableLiquidityProviderStaking()
        ImmutableFoundation()
    {}

    function initialize() public onlyFactory initializer {
        __ERC20_init("BToken", "BOB");
        _setSplitsInternal(332, 332, 332, 4);
    }

    /// @dev sets the percentage that will be divided between all the staking
    /// contracts, must only be called by _admin
    function setSplits(
        uint256 validatorStakingSplit_,
        uint256 publicStakingSplit_,
        uint256 liquidityProviderStakingSplit_,
        uint256 protocolFee_
    ) public onlyAdmin {
        _setSplitsInternal(
            validatorStakingSplit_,
            publicStakingSplit_,
            liquidityProviderStakingSplit_,
            protocolFee_
        );
    }

    /// Distributes the yields of the BToken sale to all stakeholders
    /// (miners, stakers, lp stakers, foundation, etc).
    function distribute()
        public
        returns (
            uint256 minerAmount,
            uint256 stakingAmount,
            uint256 lpStakingAmount,
            uint256 foundationAmount
        )
    {
        return _distribute();
    }

    /// Deposits a BToken amount into the AliceNet blockchain. The BTokens amount
    /// is deducted from the sender and it is burned by this contract. The
    /// created deposit Id is owned by the to_ address.
    /// @param accountType_ The type of account the to_ address must be equivalent with ( 1 for Eth native, 2 for BN )
    /// @param to_ The address of the account that will own the deposit
    /// @param amount_ The amount of BTokens to be deposited
    /// Return The deposit ID of the deposit created
    function deposit(
        uint8 accountType_,
        address to_,
        uint256 amount_
    ) public returns (uint256) {
        return _deposit(accountType_, to_, amount_);
    }

    /// Allows deposits to be minted in a virtual manner and sent to the AliceNet
    /// chain by simply emitting a Deposit event without actually minting or
    /// burning any tokens, must only be called by _admin.
    /// @param accountType_ The type of account the to_ address must be equivalent with ( 1 for Eth native, 2 for BN )
    /// @param to_ The address of the account that will own the deposit
    /// @param amount_ The amount of BTokens to be deposited
    /// Return The deposit ID of the deposit created
    function virtualMintDeposit(
        uint8 accountType_,
        address to_,
        uint256 amount_
    ) public onlyAdmin returns (uint256) {
        return _virtualDeposit(accountType_, to_, amount_);
    }

    /// Allows deposits to be minted in a virtual manner and sent to the AliceNet
    /// chain by simply emitting a Deposit event without actually minting or
    /// burning any tokens. This function receives ether in the transaction and
    /// converts them into a deposit of BToken in the AliceNet chain.
    /// This function has the same effect as calling mint (creating the
    /// tokens) + deposit (burning the tokens) functions but spending less gas.
    /// @param accountType_ The type of account the to_ address must be equivalent with ( 1 for Eth native, 2 for BN )
    /// @param to_ The address of the account that will own the deposit
    /// @param minBTK_ The amount of BTokens to be deposited
    /// Return The deposit ID of the deposit created
    function mintDeposit(
        uint8 accountType_,
        address to_,
        uint256 minBTK_
    ) public payable returns (uint256) {
        return _mintDeposit(accountType_, to_, minBTK_, msg.value);
    }

    /// Mints BToken. This function receives ether in the transaction and
    /// converts them into BToken using a bonding price curve.
    /// @param minBTK_ Minimum amount of BToken that you wish to mint given an
    /// amount of ether. If its not possible to mint the desired amount with the
    /// current price in the bonding curve, the transaction is reverted. If the
    /// minBTK_ is met, the whole amount of ether sent will be converted in BToken.
    /// Return The number of BToken minted
    function mint(uint256 minBTK_) public payable returns (uint256 numBTK) {
        numBTK = _mint(msg.sender, msg.value, minBTK_);
        return numBTK;
    }

    /// Mints BToken. This function receives ether in the transaction and
    /// converts them into BToken using a bonding price curve.
    /// @param to_ The account to where the tokens will be minted
    /// @param minBTK_ Minimum amount of BToken that you wish to mint given an
    /// amount of ether. If its not possible to mint the desired amount with the
    /// current price in the bonding curve, the transaction is reverted. If the
    /// minBTK_ is met, the whole amount of ether sent will be converted in BToken.
    /// Return The number of BToken minted
    function mintTo(address to_, uint256 minBTK_)
        public
        payable
        returns (uint256 numBTK)
    {
        numBTK = _mint(to_, msg.value, minBTK_);
        return numBTK;
    }

    /// Burn BToken. This function sends ether corresponding to the amount of
    /// BTokens being burned using a bonding price curve.
    /// @param amount_ The amount of BToken being burned
    /// @param minEth_ Minimum amount ether that you expect to receive given the
    /// amount of BToken being burned. If the amount of BToken being burned
    /// worth less than this amount the transaction is reverted.
    /// Return The number of ether being received
    function burn(uint256 amount_, uint256 minEth_)
        public
        returns (uint256 numEth)
    {
        numEth = _burn(msg.sender, msg.sender, amount_, minEth_);
        return numEth;
    }

    /// Burn BTokens and send the ether received to other account. This
    /// function sends ether corresponding to the amount of BTokens being
    /// burned using a bonding price curve.
    /// @param to_ The account to where the ether from the burning will send
    /// @param amount_ The amount of BTokens being burned
    /// @param minEth_ Minimum amount ether that you expect to receive given the
    /// amount of BTokens being burned. If the amount of BTokens being burned
    /// worth less than this amount the transaction is reverted.
    /// Return The number of ether being received
    function burnTo(
        address to_,
        uint256 amount_,
        uint256 minEth_
    ) public returns (uint256 numEth) {
        numEth = _burn(msg.sender, to_, amount_, minEth_);
        return numEth;
    }

    /// Gets the pool balance in ether
    function getPoolBalance() public view returns (uint256) {
        return _poolBalance;
    }

    /// Gets the total amount of BTokens that were deposited in the AliceNet
    /// blockchain. Since BTokens are burned when deposited, this value will be
    /// different from the total supply of BTokens.
    function getTotalBTokensDeposited() public view returns (uint256) {
        return _totalDeposited;
    }

    /// Gets the deposited amount given a depositID.
    /// @param depositID The Id of the deposit
    function getDeposit(uint256 depositID)
        public
        view
        returns (Deposit memory)
    {
        Deposit memory d = _deposits[depositID];
        require(
            d.account != address(uint160(0x00)),
            string(abi.encodePacked(BTokenErrorCodes.BTOKEN_INVALID_DEPOSIT_ID))
        );
        return d;
    }

    /// Converts an amount of BTokens in ether given a point in the bonding
    /// curve (poolbalance and totalsupply at given time).
    /// @param poolBalance_ The pool balance (in ether) at a given moment
    /// where we want to compute the amount of ether.
    /// @param totalSupply_ The total supply of BToken at a given moment
    /// where we want to compute the amount of ether.
    /// @param numBTK_ Amount of BTokens that we want to convert in ether
    function bTokensToEth(
        uint256 poolBalance_,
        uint256 totalSupply_,
        uint256 numBTK_
    ) public pure returns (uint256 numEth) {
        return _bTokensToEth(poolBalance_, totalSupply_, numBTK_);
    }

    /// Converts an amount of ether in BTokens given a point in the bonding
    /// curve (poolbalance at given time).
    /// @param poolBalance_ The pool balance (in ether) at a given moment
    /// where we want to compute the amount of BTokens.
    /// @param numEth_ Amount of ether that we want to convert in BTokens
    function ethToBTokens(uint256 poolBalance_, uint256 numEth_)
        public
        pure
        returns (uint256)
    {
        return _ethToBTokens(poolBalance_, numEth_);
    }

    /// Distributes the yields from the BToken minting to all stake holders.
    function _distribute()
        internal
        withLock
        returns (
            uint256 minerAmount,
            uint256 stakingAmount,
            uint256 lpStakingAmount,
            uint256 foundationAmount
        )
    {
        // make a local copy to save gas
        uint256 poolBalance = _poolBalance;

        // find all value in excess of what is needed in pool
        uint256 excess = address(this).balance - poolBalance;

        // take out protocolFee from excess and decrement excess
        foundationAmount = (excess * _protocolFee) / _PERCENTAGE_SCALE;

        // split remaining between miners, stakers and lp stakers
        stakingAmount = (excess * _publicStakingSplit) / _PERCENTAGE_SCALE;
        lpStakingAmount =
            (excess * _liquidityProviderStakingSplit) /
            _PERCENTAGE_SCALE;
        // then give miners the difference of the original and the sum of the
        // stakingAmount
        minerAmount =
            excess -
            (stakingAmount + lpStakingAmount + foundationAmount);

        if (foundationAmount != 0) {
            _safeTransferEthWithMagic(
                IMagicEthTransfer(_foundationAddress()),
                foundationAmount
            );
        }
        if (minerAmount != 0) {
            _safeTransferEthWithMagic(
                IMagicEthTransfer(_validatorStakingAddress()),
                minerAmount
            );
        }
        if (stakingAmount != 0) {
            _safeTransferEthWithMagic(
                IMagicEthTransfer(_publicStakingAddress()),
                stakingAmount
            );
        }
        if (lpStakingAmount != 0) {
            _safeTransferEthWithMagic(
                IMagicEthTransfer(_liquidityProviderStakingAddress()),
                lpStakingAmount
            );
        }
        require(
            address(this).balance >= poolBalance,
            string(abi.encodePacked(BTokenErrorCodes.BTOKEN_INVALID_BALANCE))
        );

        // invariants hold
        return (minerAmount, stakingAmount, lpStakingAmount, foundationAmount);
    }

    // Burn the tokens during deposits without sending ether back to user as the
    // normal burn function. The ether will be distributed in the distribute
    // method.
    function _destroyTokens(uint256 numBTK_) internal returns (bool) {
        require(
            numBTK_ != 0,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_INVALID_BURN_AMOUNT)
            )
        );
        _poolBalance -= _bTokensToEth(_poolBalance, totalSupply(), numBTK_);
        ERC20Upgradeable._burn(msg.sender, numBTK_);
        return true;
    }

    // Internal function that does the deposit in the AliceNet Chain, i.e emit the
    // event DepositReceived. All the BTokens sent to this function are burned.
    function _deposit(
        uint8 accountType_,
        address to_,
        uint256 amount_
    ) internal returns (uint256) {
        require(
            !_isContract(to_),
            string(
                abi.encodePacked(
                    BTokenErrorCodes.BTOKEN_CONTRACTS_DISALLOWED_DEPOSITS
                )
            )
        );
        require(
            amount_ > 0,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_DEPOSIT_AMOUNT_ZERO)
            )
        );
        require(
            _destroyTokens(amount_),
            string(abi.encodePacked(BTokenErrorCodes.BTOKEN_DEPOSIT_BURN_FAIL))
        );
        // copying state to save gas
        return _doDepositCommon(accountType_, to_, amount_);
    }

    // does a virtual deposit into the AliceNet Chain without actually minting or
    // burning any token.
    function _virtualDeposit(
        uint8 accountType_,
        address to_,
        uint256 amount_
    ) internal returns (uint256) {
        require(
            !_isContract(to_),
            string(
                abi.encodePacked(
                    BTokenErrorCodes.BTOKEN_CONTRACTS_DISALLOWED_DEPOSITS
                )
            )
        );
        require(
            amount_ > 0,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_DEPOSIT_AMOUNT_ZERO)
            )
        );
        // copying state to save gas
        return _doDepositCommon(accountType_, to_, amount_);
    }

    // Mints a virtual deposit into the AliceNet Chain without actually minting or
    // burning any token. This function converts ether sent in BTokens.
    function _mintDeposit(
        uint8 accountType_,
        address to_,
        uint256 minBTK_,
        uint256 numEth_
    ) internal returns (uint256) {
        require(
            !_isContract(to_),
            string(
                abi.encodePacked(
                    BTokenErrorCodes.BTOKEN_CONTRACTS_DISALLOWED_DEPOSITS
                )
            )
        );
        require(
            numEth_ >= _MARKET_SPREAD,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_MARKET_SPREAD_TOO_LOW)
            )
        );
        numEth_ = numEth_ / _MARKET_SPREAD;
        uint256 amount_ = _ethToBTokens(_poolBalance, numEth_);
        require(
            amount_ >= minBTK_,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_MINT_INSUFFICIENT_ETH)
            )
        );
        return _doDepositCommon(accountType_, to_, amount_);
    }

    function _doDepositCommon(
        uint8 accountType_,
        address to_,
        uint256 amount_
    ) internal returns (uint256) {
        uint256 depositID = _depositID + 1;
        _deposits[depositID] = _newDeposit(accountType_, to_, amount_);
        _totalDeposited += amount_;
        _depositID = depositID;
        emit DepositReceived(depositID, accountType_, to_, amount_);
        return depositID;
    }

    // Internal function that mints the BToken tokens following the bounding
    // price curve.
    function _mint(
        address to_,
        uint256 numEth_,
        uint256 minBTK_
    ) internal returns (uint256 numBTK) {
        require(
            numEth_ >= _MARKET_SPREAD,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_MARKET_SPREAD_TOO_LOW)
            )
        );
        numEth_ = numEth_ / _MARKET_SPREAD;
        uint256 poolBalance = _poolBalance;
        numBTK = _ethToBTokens(poolBalance, numEth_);
        require(
            numBTK >= minBTK_,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_MINIMUM_MINT_NOT_MET)
            )
        );
        poolBalance += numEth_;
        _poolBalance = poolBalance;
        ERC20Upgradeable._mint(to_, numBTK);
        return numBTK;
    }

    // Internal function that burns the BToken tokens following the bounding
    // price curve.
    function _burn(
        address from_,
        address to_,
        uint256 numBTK_,
        uint256 minEth_
    ) internal returns (uint256 numEth) {
        require(
            numBTK_ != 0,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_INVALID_BURN_AMOUNT)
            )
        );
        uint256 poolBalance = _poolBalance;
        numEth = _bTokensToEth(poolBalance, totalSupply(), numBTK_);
        require(
            numEth >= minEth_,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_MINIMUM_BURN_NOT_MET)
            )
        );
        poolBalance -= numEth;
        _poolBalance = poolBalance;
        ERC20Upgradeable._burn(from_, numBTK_);
        _safeTransferEth(to_, numEth);
        return numEth;
    }

    function _setSplitsInternal(
        uint256 validatorStakingSplit_,
        uint256 publicStakingSplit_,
        uint256 liquidityProviderStakingSplit_,
        uint256 protocolFee_
    ) internal {
        require(
            validatorStakingSplit_ +
                publicStakingSplit_ +
                liquidityProviderStakingSplit_ +
                protocolFee_ ==
                _PERCENTAGE_SCALE,
            string(
                abi.encodePacked(BTokenErrorCodes.BTOKEN_SPLIT_VALUE_SUM_ERROR)
            )
        );
        _validatorStakingSplit = validatorStakingSplit_;
        _publicStakingSplit = publicStakingSplit_;
        _liquidityProviderStakingSplit = liquidityProviderStakingSplit_;
        _protocolFee = protocolFee_;
    }

    // Check if addr_ is EOA (Externally Owned Account) or a contract.
    function _isContract(address addr_) internal view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(addr_)
        }
        return size > 0;
    }

    // Internal function that converts an ether amount into BToken tokens
    // following the bounding price curve.
    function _ethToBTokens(uint256 poolBalance_, uint256 numEth_)
        internal
        pure
        returns (uint256)
    {
        return _fx(poolBalance_ + numEth_) - _fx(poolBalance_);
    }

    // Internal function that converts a BToken amount into ether following the
    // bounding price curve.
    function _bTokensToEth(
        uint256 poolBalance_,
        uint256 totalSupply_,
        uint256 numBTK_
    ) internal pure returns (uint256 numEth) {
        require(
            totalSupply_ >= numBTK_,
            string(
                abi.encodePacked(
                    BTokenErrorCodes.BTOKEN_BURN_AMOUNT_EXCEEDS_SUPPLY
                )
            )
        );
        return
            _min(poolBalance_, _fp(totalSupply_) - _fp(totalSupply_ - numBTK_));
    }

    function _newDeposit(
        uint8 accountType_,
        address account_,
        uint256 value_
    ) internal pure returns (Deposit memory) {
        Deposit memory d = Deposit(accountType_, account_, value_);
        return d;
    }
}
