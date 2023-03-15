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
| read_transaction          | :white_check_mark: |
| transaction_size          | :white_check_mark: |
| expiration                | :white_check_mark: |
| tapos_block_num           | :white_check_mark: |
| tapos_block_prefix        | :white_check_mark: |
| get_action                | :white_check_mark: |

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
| get_context_free_data     | :white_check_mark: |

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
| db_idx64_store            | :white_check_mark: |
| db_idx64_update           | :white_check_mark: |
| db_idx64_remove           | :white_check_mark: |
| db_idx64_find_secondary   | :white_check_mark: |
| db_idx64_find_primary     | :white_check_mark: |
| db_idx64_lowerbound       | :white_check_mark: |
| db_idx64_upperbound       | :white_check_mark: |
| db_idx64_end              | :white_check_mark: |
| db_idx64_next             | :white_check_mark: |
| db_idx64_previous         | :white_check_mark: |
| db_idx128_store           | :white_check_mark: |
| db_idx128_update          | :white_check_mark: |
| db_idx128_remove          | :white_check_mark: |
| db_idx128_find_secondary  | :white_check_mark: |
| db_idx128_find_primary    | :white_check_mark: |
| db_idx128_lowerbound      | :white_check_mark: |
| db_idx128_upperbound      | :white_check_mark: |
| db_idx128_end             | :white_check_mark: |
| db_idx128_next            | :white_check_mark: |
| db_idx128_previous        | :white_check_mark: |
| db_idx256_store           | :white_check_mark: |
| db_idx256_update          | :white_check_mark: |
| db_idx256_remove          | :white_check_mark: |
| db_idx256_find_secondary  | :white_check_mark: |
| db_idx256_find_primary    | :white_check_mark: |
| db_idx256_lowerbound      | :white_check_mark: |
| db_idx256_upperbound      | :white_check_mark: |
| db_idx256_end             | :white_check_mark: |
| db_idx256_next            | :white_check_mark: |
| db_idx256_previous        | :white_check_mark: |
| db_idx_double_store                | :white_check_mark: |
| db_idx_double_update                | :white_check_mark: |
| db_idx_double_remove                | :white_check_mark: |
| db_idx_double_find_secondary                | :white_check_mark: |
| db_idx_double_find_primary                | :white_check_mark: |
| db_idx_double_lowerbound                | :white_check_mark: |
| db_idx_double_upperbound                | :white_check_mark: |
| db_idx_double_end                | :white_check_mark: |
| db_idx_double_next                | :white_check_mark: |
| db_idx_double_previous                | :white_check_mark: |
| db_idx_long_double_store                | :white_check_mark: |
| db_idx_long_double_update                | :white_check_mark: |
| db_idx_long_double_remove                | :white_check_mark: |
| db_idx_long_double_find_secondary                | :white_check_mark: |
| db_idx_long_double_find_primary                | :white_check_mark: |
| db_idx_long_double_lowerbound                | :white_check_mark: |
| db_idx_long_double_upperbound                | :white_check_mark: |
| db_idx_long_double_end                | :white_check_mark: |
| db_idx_long_double_next                | :white_check_mark: |
| db_idx_long_double_previous                | :white_check_mark: |

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
| is_feature_active               | :white_check_mark: |
| activate_feature               | :white_check_mark: |
| preactivate_feature               | Missing |
| set_resource_limits               | :white_check_mark: |
| get_resource_limits            | :white_check_mark: |
| get_wasm_parameters_packed                     | Missing |
| set_wasm_parameters_packed                    | Missing |
| set_proposed_producers  | Missing |
| set_proposed_producers_ex  | Missing |
| get_blockchain_parameters_packed          | Missing |
| set_blockchain_parameters_packed               | Missing |
| get_parameters_packed               | Missing |
| set_parameters_packed               | Missing |
| is_privileged               | :white_check_mark: |
| set_privileged               | :white_check_mark: |

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
| send_inline               | :white_check_mark: |
| send_context_free_inline  | :white_check_mark: |
| send_deferred             | Missing |
| cancel_deferred           | Missing |

### Compiler builtins:

| Name                      | Status      |
| ------------------------- |:-----------:|
| __ashlti3                 | :white_check_mark: |
| __ashrti3                 | :white_check_mark: |
| __lshlti3                 | :white_check_mark: |
| __lshrti3                 | :white_check_mark: |
| __divti3                  | :white_check_mark: |
| __udivti3                 | :white_check_mark: |
| __multi3                  | :white_check_mark: |
| __modti3                  | :white_check_mark: |
| __umodti3                 | :white_check_mark: |
| __addtf3                  | :white_check_mark: |
| __subtf3                  | :white_check_mark: |
| __multf3                  | :white_check_mark: |
| __divtf3                  | :white_check_mark: |
| __negtf2                  | :white_check_mark: |
| __extendsftf2             | :white_check_mark: |
| __extenddftf2             | :white_check_mark: |
| __trunctfdf2              | :white_check_mark: |
| __trunctfsf2              | :white_check_mark: |
| __fixtfsi                 | :white_check_mark: |
| __fixtfdi                 | :white_check_mark: |
| __fixtfti                 | :white_check_mark: |
| __fixunstfsi              | :white_check_mark: |
| __fixunstfdi              | :white_check_mark: |
| __fixunstfti              | :white_check_mark: |
| __fixsfti                 | :white_check_mark: |
| __fixdfti                 | :white_check_mark: |
| __fixunssfti              | :white_check_mark: |
| __fixunsdfti              | :white_check_mark: |
| __floatsidf               | :white_check_mark: |
| __floatsitf               | :white_check_mark: |
| __floatditf               | :white_check_mark: |
| __floatunsitf             | :white_check_mark: |
| __floatunditf             | :white_check_mark: |
| __floattidf               | :white_check_mark: |
| __floatuntidf             | :white_check_mark: |
| __eqtf2                   | :white_check_mark: |
| __netf2                   | :white_check_mark: |
| __getf2                   | :white_check_mark: |
| __gttf2                   | :white_check_mark: |
| __letf2                   | :white_check_mark: |
| __lttf2                   | :white_check_mark: |
| __cmptf2                  | :white_check_mark: |
| __unordtf2                | Missing |