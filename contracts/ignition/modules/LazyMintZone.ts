// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
import LazyMint1155Module from "./LazyMint1155";

const LazyMintZoneModule = buildModule("lazyMintZone", (m) => {
  // // Import the LazyMint1155 module to reference its deployed contract
  const lazyMint1155ModuleInstance = m.useModule(LazyMint1155Module);

  // // Deploy the LazyMintZone contract, passing the LazyMint1155 contract address

  const ownerAddress = m.getParameter(
    "signerAddress",
    "0x253F9Dd15f4Bd360595b0E83d51ef31d8E71d31B"
  );
  const lazyMintZoneV1 = m.contract("LazyMintZone", [
    ownerAddress,
    lazyMint1155ModuleInstance.lazyMint1155,
  ]);

  return {
    lazyMintZoneV1,
    lazyMint1155: lazyMint1155ModuleInstance.lazyMint1155,
  };
});

export default LazyMintZoneModule;
