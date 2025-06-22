// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const LazyMint1155Module = buildModule("LazyMint1155ModuleV1", (m) => {
  const signerAddress = m.getParameter(
    "signerAddress",
    "0x253F9Dd15f4Bd360595b0E83d51ef31d8E71d31B"
  );

  const lazyMint1155 = m.contract("LazyMint1155", [signerAddress]);

  return { lazyMint1155 };
});

export default LazyMint1155Module;
