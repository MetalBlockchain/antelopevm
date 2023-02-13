# Antelope VM

Antelope based Virtual Machine for the Metal Blockchain to support the A chain. At its core it will be capable of running Antelope / Proton transactions against WebAssembly based smart-contracts.

**This is work in progress**

## Database format

Antelope VM relies on BadgerDB as its key-value store having access to the entire DB in-memory.

## Implemented Host Functions

### Action functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| read_action_data          | :white_check_mark: |
| action_data_size          | :white_check_mark: |
| current_receiver          | :white_check_mark: |
| set_action_return_value   | :white_check_mark: |

### Authorization functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| require_auth              | :white_check_mark: |
| has_auth                  | :white_check_mark: |
| require_auth2             | :white_check_mark: |
| require_recipient         | :white_check_mark: |
| is_account                | :white_check_mark: |

### Context-free system functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| abort                     | :white_check_mark: |
| eosio_assert              | :white_check_mark: |
| eosio_assert_message      | :white_check_mark: |
| eosio_assert_code         | :white_check_mark: |
| eosio_exit                | :white_check_mark: |

### Context-free transaction functions:

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
| prints                    | :white_check_mark: |
| prints_l                  | :white_check_mark: |
| printi                    | :white_check_mark: |
| printui                   | :white_check_mark: |
| printi128                 | :white_check_mark: |
| printui128                | :white_check_mark: |
| printsf                   | :white_check_mark: |
| printdf                   | :white_check_mark: |
| printqf                   | :white_check_mark: |
| printn                    | :white_check_mark: |
| printhex                  | :white_check_mark: |

### Context-free functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| get_context_free_data     | Missing |

### Crypto functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| assert_recover_key        | :white_check_mark: |
| recover_key               | :white_check_mark: |
| assert_sha256             | :white_check_mark: |
| assert_sha1               | :white_check_mark: |
| assert_sha512             | :white_check_mark: |
| assert_ripemd160          | :white_check_mark: |
| sha1                      | :white_check_mark: |
| sha256                    | :white_check_mark: |
| sha512                    | :white_check_mark: |
| ripemd160                 | :white_check_mark: |
| alt_bn128_add                 | Missing |
| alt_bn128_mul                 | Missing |
| alt_bn128_pair                 | Missing |
| mod_exp                 | Missing |
| blake2_f                 | Missing |
| sha3                 | Missing |
| k1_recover                 | Missing |

### Database functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| db_store_i64              | :white_check_mark: |
| db_update_i64             | :white_check_mark: |
| db_remove_i64             | :white_check_mark: |
| db_get_i64                | :white_check_mark: |
| db_next_i64               | :white_check_mark: |
| db_previous_i64           | :white_check_mark: |
| db_find_i64               | :white_check_mark: |
| db_lowerbound_i64         | :white_check_mark: |
| db_upperbound_i64         | :white_check_mark: |
| db_idx64_store            | Missing |
| db_idx64_update           | Missing |
| db_idx64_remove           | Missing |
| db_idx64_find_secondary   | Missing |
| db_idx64_find_primary     | Missing |
| db_idx64_lowerbound                | Missing |
| db_idx64_upperbound                | Missing |
| db_idx64_end                | Missing |
| db_idx64_next                | Missing |
| db_idx64_previous                | Missing |
| db_idx128_store           | Missing |
| db_idx128_update          | Missing |
| db_idx128_remove          | Missing |
| db_idx128_find_secondary  | Missing |
| db_idx128_find_primary                | Missing |
| db_idx128_lowerbound                | Missing |
| db_idx128_upperbound                | Missing |
| db_idx128_end                | Missing |
| db_idx128_next                | Missing |
| db_idx128_previous                | Missing |
| db_idx256_store                | Missing |
| db_idx256_update                | Missing |
| db_idx256_remove                | Missing |
| db_idx256_find_secondary                | Missing |
| db_idx256_find_primary                | Missing |
| db_idx256_lowerbound                | Missing |
| db_idx256_upperbound                | Missing |
| db_idx256_end                | Missing |
| db_idx256_next                | Missing |
| db_idx256_previous                | Missing |
| db_idx_double_store                | Missing |
| db_idx_double_update                | Missing |
| db_idx_double_remove                | Missing |
| db_idx_double_find_secondary                | Missing |
| db_idx_double_find_primary                | Missing |
| db_idx_double_lowerbound                | Missing |
| db_idx_double_upperbound                | Missing |
| db_idx_double_end                | Missing |
| db_idx_double_next                | Missing |
| db_idx_double_previous                | Missing |
| db_idx_long_double_store                | Missing |
| db_idx_long_double_update                | Missing |
| db_idx_long_double_remove                | Missing |
| db_idx_long_double_find_secondary                | Missing |
| db_idx_long_double_find_primary                | Missing |
| db_idx_long_double_lowerbound                | Missing |
| db_idx_long_double_upperbound                | Missing |
| db_idx_long_double_end                | Missing |
| db_idx_long_double_next                | Missing |
| db_idx_long_double_previous                | Missing |

### Memory functions:

| Name                      | Status      |
| ------------------------- |:-----------:|
| memcpy                    | :white_check_mark: |
| memmove                   | :white_check_mark: |
| memcmp                    | :white_check_mark: |
| memset                    | :white_check_mark: |

### Permission functions:

| Name                              | Status  |
| --------------------------------- |:-------:|
| check_transaction_authorization   | :white_check_mark: |
| check_permission_authorization    | :white_check_mark: |
| get_permission_last_used          | :white_check_mark: |
| get_account_creation_time         | :white_check_mark: |

### Privileged functions:

| Name                              | Status  |
| --------------------------------- |:-------:|
| is_feature_active               | Missing |
| preactivate_feature               | Missing |
| set_resource_limits               | Missing |
| get_resource_limits            | Missing |
| set_proposed_producers_common         | Missing |
| get_wasm_parameters_packed                     | Missing |
| set_wasm_parameters_packed                    | Missing |
| set_proposed_producers  | Missing |
| set_proposed_producers_ex  | Missing |
| get_blockchain_parameters_packed          | Missing |
| set_blockchain_parameters_packed               | Missing |
| get_parameters_packed               | Missing |
| set_parameters_packed               | Missing |
| is_privileged               | Missing |
| set_privileged               | Missing |

### Producer functions:

| Name                              | Status  |
| --------------------------------- |:-------:|
| get_active_producers              | Missing |

### System functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| current_time              | :white_check_mark: |
| publication_time          | Missing |
| is_feature_activated      | :white_check_mark: |
| get_sender                | :white_check_mark: |
| get_block_num             | Missing |

### Transaction functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| send_inline               | Missing |
| send_context_free_inline  | Missing |
| send_deferred             | Missing |
| cancel_deferred           | Missing |

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