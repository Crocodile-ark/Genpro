# GXR (Gen X Raider) Blockchain Specification — FINAL

> ✅ Designed without smart contracts, inflation-proof, based on PoS & IBC
> ✅ Focused on efficiency, fair distribution, and automated decentralization

---

## 1. CORE IDENTIFICATION

| Component            | Details                                   |
| -------------------- | ------------------------------------------ |
| **Chain Name**       | Gen X Raider (GXR)                         |
| **Token Ticker**     | GXR                                        |
| **Smallest Denom**   | `gen` (1 GXR = 100,000,000 gen)            |
| **Total Supply**     | 85,000,000 GXR *(fixed, non-inflationary)*|
| **Decimals**         | 8                                          |
| **Smart Contracts**  | Not used                                  |
| **Consensus**        | Proof-of-Stake (PoS)                       |
| **IBC Support**      | Active (GXR/TON, GXR/ATOM, etc.)       |
| **Max Validators**   | 85 nodes                                   |
| **Block Time**       | 15 seconds per block                      |

---

## 2. GXR TOKENOMICS

### Total Supply: 85,000,000 GXR

| Allocation                     | Amount (GXR) | Percentage | Description                                                                      |
| ------------------------------ | -------------| -----------| --------------------------------------------------------------------------------- |
| Airdrop & Farming              | 17,000,000    | 20.00%     | Initial distribution via Telegram bot farming                                     |
| Developer                     | 8,500,000     | 10.00%     | 5-year hard vesting, 10% unlocked every 6 months                                 |
| Core Team (2 members)         | 3,400,000     | 4.00%      | 2% / 2%, 3-year soft vesting                                                      |
| LP & Market                   | 8,500,000     | 10.00%     | Initial liquidity (GXR/TON, GXR/ATOM, etc.)                                   |
| Grants (5–10 parties)         | 8,500,000     | 10.00%     | Project and partner collaboration grants                                          |
| Staking Pool (PoS)            | 8,500,000     | 10.00%     | Rewards for active delegators                                                     |
| Halving Fund                  | 21,250,000    | 25.00%     | Reward pool for 5-year halving cycles                                             |
| Reserve/Expansion             | 8,500,000     | 10.00%     | Emergency and ecosystem development funds                                         |
| Early Validators (30)         | 850,000       | 1.00%      | 0.5% in year 1 and 0.5% in year 2 (if >20 active days/month)                      |
| **Multisig Mainnet Fund**     | 850,000       | 1.00%      | Unlocked, managed by multisig wallet for post-mainnet operational funding         |

---

## 3. DYNAMIC HALVING SYSTEM

### Halving Cycle: Every 5 Years

- All rewards are distributed from the Halving Fund (21,250,000 GXR)
- No new minting (fixed supply)

### Gradual Reduction Mechanism:

| Halving | Period      | Fund (GXR) | Reduction (%) | Monthly (for 24 months) |
| ------- | ----------- | -----------|----------------|--------------------------|
| 1       | Year 1–5    | 4,250,000  | —              | 177,083                  |
| 2       | Year 6–10   | 3,612,500  | -15%           | 150,521                  |
| 3       | Year 11–15  | 3,070,625  | -15%           | 127,943                  |
| 4       | Year 16–20  | 2,610,032  | -15%           | 108,751                  |
| 5       | Year 21–25  | 2,218,528  | -15%           | 92,439                   |

### Distribution Rules per Halving Cycle:

- Rewards are distributed for the **first 24 months** only in each halving cycle
- Remaining tokens are held until the next halving cycle
- Objective: Price stability, anti-dump, and strategic pause

- 70% → Active Validators
- 20% → PoS Pool (delegators)
- 10% → DEX Pool (GXR/TON, GXR/ATOM, etc.)
  - Year 1: 5% distributed evenly
  - Year 2: 5% dynamically based on volume
  - Year 3–5: no distribution (held for next halving cycle)

---

## 4. VALIDATORS & DELEGATORS

### Validators:

- Max Nodes: 85
- Initial Commission: 5%–10% (dynamic)
- Rewards: Monthly from halving + transaction fees
- Unstake Fee from Delegators: 0.5% (goes to validator)
- Inactive >10 days/month → no reward for that month

### Delegators:

- Earn rewards from PoS Pool
- Choose from active validators

---

## 5. TRANSACTION FEE SYSTEM

### A. General Transaction Fees

Used for standard transactions:

| Component              | Percentage |
| ---------------------- | -----------|
| Validator              | 40%        |
| DEX Pool (Auto Refill) | 30%        |
| PoS Pool (Delegators)  | 30%        |

### B. Community LP Farming Fees

For user-created liquidity pools:

| Component              | Percentage |
| ---------------------- | -----------|
| Validator              | 30%        |
| DEX Pool (Auto Refill) | 25%        |
| Community LP (reward)  | 25%        |
| PoS Pool (Delegators)  | 20%        |

---

## 6. LP SYSTEM: TEAM & COMMUNITY

### A. Official LP

- Created by GXR team
- Initial liquidity from 10% LP & Market allocation
- Managed by validator bot (auto refill, auto balancing)
- Uses default fee split (40/30/30)

### B. Community LPs

- Anyone can create GXR/TOKEN LPs
- Incentives based on farming fee model (30/25/25/20)
- Detected by validator bots and rewarded automatically
- No smart contract needed, only whitelisted LP address

---

## 7. VALIDATOR BOT (MANDATORY)

### Functions:

- Auto IBC Relayer (cross-chain sync)
- Auto price rebalancing across pools
- Auto reward distribution
- Auto refill DEX Pool from fee income
- Telegram alerts for uptime, pool imbalance, etc.
- One-click node deployment + bot auto-activation

> Blockchain is immutable, but bot rules can be updated by the core team.

### Bot Protection:

- Max daily swap cap (e.g., 10,000 GXR/day)
- 30-min cooldown if extreme price spikes detected
- If GXR price exceeds $5–$10 → **Bot goes into 24h monitor-only mode**
- After cooldown, bot resumes → performs rebalancing
- Reward distribution continues unaffected
- Goal: prevent price manipulation & protect system during extreme markets

---

## 8. GXR KEY FEATURES

| Feature                 | Description                                                                                                     |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------- |
| No Smart Contract       | Simple, lightweight, secure                                                                                      |
| Fixed Supply            | Inflation-proof                                                                                                   |
| Auto Fee Refill         | DEX Pools automatically refilled                                                                                 |
| Decentralized Bot       | Node + bot auto-activates on deploy, no partial options (must run full validator + bot)                         |
| Anti-Dump Halving       | Rewards spread over 2 years, not dumped instantly                                                                |
| Early Validator Bonus   | 1% bonus distributed to 30 early nodes over 2 years (if uptime >20 days/month)                                  |
| Community LPs           | Community-driven LPs with auto fee farming rewards                                                               |
| **Immutable Genesis**   | All parameters locked at launch, no SC/governance upgrades, only hard fork for emergency changes                |

