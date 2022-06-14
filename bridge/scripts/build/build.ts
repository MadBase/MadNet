import { exec } from "child_process";
import { keccak256 } from "ethers/lib/utils";
import fs from "fs/promises";
import path from "path";

const bindingsDir = path.resolve(__dirname, "../../bindings");
const abiDir = path.resolve(bindingsDir, "bindings-artifacts");
const mappingDir = path.resolve(bindingsDir, "mapping");

type signatureAssociation = { signature: Buffer; name: string };

async function generateMappingFile() {
  console.log("generating signature mapping...");

  const files = (await fs.readdir(abiDir)).filter((file) =>
    file.endsWith(".json")
  );

  const associations: signatureAssociation[] = (
    await Promise.all(
      files.map(async (fileName) => {
        const contents = await fs.readFile(path.resolve(abiDir, fileName));

        const abi = JSON.parse(contents.toString());
        const contractName = path.basename(fileName, ".json");

        return abi
          .filter((x: any) => x.type === "function")
          .map((func: any) => {
            const canonical = `${func.name}(${func.inputs
              .map((input: any) => input.type)
              .join(",")})`;

            return {
              signature: Buffer.from(
                keccak256(Buffer.from(canonical)).slice(2, 10),
                "hex"
              ),
              name: `${contractName}.${canonical}`,
            };
          });
      })
    )
  ).flat();

  const golangSrc = `// Generated by custom bridge script. DO NOT EDIT.

package signatures

var FunctionMapping = map[[4]byte]string {
${associations
  .map(({ name, signature }) => `\t{${signature.join(",")}}: "${name}"`)
  .join(",\n")}}
`;

  await fs.mkdir(mappingDir, { recursive: true });
  await fs.writeFile(
    path.resolve(mappingDir, "function_mapping.go"),
    golangSrc
  );
}

async function generateIfaces() {
  const contracts = (await fs.readdir(abiDir))
    .filter((file) => file.endsWith(".json"))
    .map((file) => path.basename(file, ".json"));

  for (const contract of contracts) {
    console.log(`generating ${contract} interface...`);
    await new Promise((resolve, reject) =>
      exec(
        `ifacemaker -f ${contract}.go -s ${contract}Caller -o ${contract}.iface.c.go -p bindings -i I${contract}Caller -c "Generated by ifacemaker. DO NOT EDIT." &&
         ifacemaker -f ${contract}.go -s ${contract}Transactor -o ${contract}.iface.t.go -p bindings -i I${contract}Transactor -c "Generated by ifacemaker. DO NOT EDIT." &&
         ifacemaker -f ${contract}.go -s ${contract}Filterer -o ${contract}.iface.f.go -p bindings -i I${contract}Filterer -c "Generated by ifacemaker. DO NOT EDIT."`,
        { cwd: bindingsDir },
        (error, stdout, stderr) => {
          if (error) reject(Error(stderr));
          else resolve(stdout);
        }
      )
    );
  }

  const golangSrc = `// Generated by custom bridge script. DO NOT EDIT

package bindings
  
${contracts
  .map(
    (contract) => `type I${contract} interface {
\tI${contract}Caller
\tI${contract}Transactor
\tI${contract}Filterer
}`
  )
  .join("\n\n")}
`;

  await fs.writeFile(path.resolve(bindingsDir, "iface.go"), golangSrc);
}

async function main() {
  await generateMappingFile();
  await generateIfaces();
}
void main();
