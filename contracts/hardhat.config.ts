import "dotenv/config";
import { HardhatUserConfig, vars } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "@nomicfoundation/hardhat-verify";
import "solidity-coverage";

const ETHERSCAN_API_KEY = vars.get("ETHERSCAN_API_KEY");

const config: HardhatUserConfig = {
  solidity: "0.8.28",
  etherscan: {
    apiKey: ETHERSCAN_API_KEY,
  },

  typechain: {
    target: "ethers-v6",
  },
  networks: {
    hardhat: {
      chainId: 1337, // Default Hardhat network chain ID
      gasPrice: 20000000000, // 20 Gwei

      accounts: [
        {
          privateKey: vars.get("DEPLOYER_PRIVATE_KEY"),
          balance: "1000000000000000000000", // 1000 ETH
        },
        {
          privateKey: vars.get("SIGNER_PRIVATE_KEY"),
          balance: "1000000000000000000000", // 1000 ETH
        },
        {
          privateKey: vars.get("ACCOUNT1_PRIVATE_KEY"),
          balance: "1000000000000000000000", // 1000 ETH
        },
      ], // Uncomment to use private key directly
    },

    bscTestnet: {
      url: "https://data-seed-prebsc-1-s1.binance.org:8545", // BSC testnet endpoint
      chainId: 97,
      gasPrice: 20000000000,
      accounts: [vars.get("DEPLOYER_PRIVATE_KEY")], // Use .env for private key
    },
  },
};

export default config;
