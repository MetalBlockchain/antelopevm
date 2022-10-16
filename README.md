# Antelope VM

Antelope based Virtual Machine for the Metal Blockchain to support the A chain. At its core it will be capable of running Antelope / Proton transactions against WebAssembly based smart-contracts.

**This is work in progress**

## Implemented Host Functions

### Action functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| read_action_data          | Missing |
| action_data_size          | Missing |
| current_receiver          | Missing |
| set_action_return_value   | Missing |

### Authorization functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| require_auth              | Missing |
| has_auth                  | Missing |
| require_auth2             | Missing |
| require_recipient         | Missing |
| is_account                | Missing |

### Assert functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| eosio_assert              | Missing |
| eosio_assert_message      | Missing |
| eosio_assert_code         | Missing |
| eosio_exit                | Missing |

### Transaction functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| read_transaction          | Missing |
| transaction_size          | Missing |
| expiration                | Missing |
| tapos_block_num           | Missing |
| tapos_block_prefix        | Missing |
| get_action                | Missing |

### Console functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| prints                    | Missing |
| prints_l                  | Missing |
| printi                    | Missing |
| printui                   | Missing |
| printi128                 | Missing |
| printui128                | Missing |
| printsf                   | Missing |
| printdf                   | Missing |
| printqf                   | Missing |
| printn                    | Missing |
| printhex                  | Missing |

### Context free functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| get_context_free_data     | Missing |

### Crypto functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| assert_recover_key        | Missing |
| recover_key               | Missing |
| assert_sha256             | Missing |
| assert_sha1               | Missing |
| assert_sha512             | Missing |
| assert_ripemd160          | Missing |
| sha1                      | Missing |
| sha256                    | Missing |
| sha512                    | Missing |
| ripemd160                 | Missing |

### Compiler builtins:

| Name                      | Status      |
| ------------------------- |:-----------:|
| __ashlti3                 | Implemented |
| __ashrti3                 | Implemented |
| __lshlti3                 | Implemented |
| __lshrti3                 | Implemented |
| __divti3                  | Implemented |
| __udivti3                 | Implemented |
| __multi3                  | Implemented |
| __modti3                  | Implemented |
| __umodti3                 | Implemented |
| __addtf3                  | Implemented |
| __subtf3                  | Implemented |
| __multf3                  | Implemented |
| __divtf3                  | Implemented |
| __negtf2                  | Implemented |
| __extendsftf2             |   Missing  |
| __extenddftf2                 | Missing |
| __trunctfdf2                 | Missing |
| __trunctfsf2                 | Missing |
| __fixtfsi                 | Missing |
| __fixtfdi                 | Missing |
| __fixtfti                 | Missing |
| __fixunstfsi                 | Missing |
| __fixunstfdi                 | Missing |
| __fixunstfti                 | Missing |
| __fixsfti                 | Missing |
| __fixdfti                 | Missing |
| __fixunssfti                 | Missing |
| __fixunsdfti                 | Missing |
| __floatsidf                 | Missing |
| __floatsitf                 | Missing |
| __floatditf                 | Missing |
| __floatunsitf                 | Missing |
| __floatunditf                 | Missing |
| __floattidf                 | Missing |
| __floatuntidf                 | Missing |
| __eqtf2                 | Missing |
| __netf2                 | Missing |
| __getf2                 | Missing |
| __gttf2                 | Missing |
| __letf2                 | Missing |
| __lttf2                 | Missing |
| __cmptf2                 | Missing |
| __unordtf2                 | Missing |