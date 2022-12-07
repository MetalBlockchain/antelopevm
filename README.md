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

### Database functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| db_store_i64              | Missing |
| db_update_i64             | Missing |
| db_remove_i64             | Missing |
| db_get_i64                | Missing |
| db_next_i64               | Missing |
| db_previous_i64           | Missing |
| db_find_i64               | Missing |
| db_lowerbound_i64         | Missing |
| db_upperbound_i64         | Missing |
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

### Key value functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| kv_erase              | Missing |
| kv_set          | Missing |
| kv_get      | Missing |
| kv_get_data                | Missing |
| kv_it_create                | Missing |
| kv_it_destroy                | Missing |
| kv_it_status                | Missing |
| kv_it_compare                | Missing |
| kv_it_key_compare                | Missing |
| kv_it_move_to_end                | Missing |
| kv_it_next                | Missing |
| kv_it_prev                | Missing |
| kv_it_lower_bound                | Missing |
| kv_it_key                | Missing |
| kv_it_value                | Missing |

### Memory functions:

| Name                      | Status      |
| ------------------------- |:-----------:|
| memcpy                    | Implemented |
| memmove                   | Implemented |
| memcmp                    | Implemented |
| memset                    | Implemented |

### Permission functions:

| Name                              | Status  |
| --------------------------------- |:-------:|
| check_transaction_authorization   | Missing |
| check_permission_authorization    | Missing |
| get_permission_last_used          | Missing |
| get_account_creation_time         | Missing |

### Privileged functions:

| Name                              | Status  |
| --------------------------------- |:-------:|
| is_feature_active                 | Missing |
| preactivate_feature               | Missing |
| set_resource_limits               | Missing |
| get_resource_limits               | Missing |
| set_resource_limit                | Missing |
| get_resource_limit                | Missing |
| get_wasm_parameters_packed        | Missing |
| set_wasm_parameters_packed        | Missing |
| set_proposed_producers            | Missing |
| set_proposed_producers_ex         | Missing |
| get_blockchain_parameters_packed  | Missing |
| set_blockchain_parameters_packed  | Missing |
| get_parameters_packed             | Missing |
| set_parameters_packed             | Missing |
| get_kv_parameters_packed          | Missing |
| set_kv_parameters_packed          | Missing |
| is_privileged                     | Missing |
| set_privileged                    | Missing |

### Producer functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| get_active_producers      | Missing |

### System functions:

| Name                      | Status  |
| ------------------------- |:-------:|
| current_time              | Missing |
| publication_time          | Missing |
| is_feature_activated      | Missing |
| get_sender                | Missing |

### Transaction functions:

| Name                      | Status                    |
| ------------------------- |:-------------------------:|
| send_inline               | Missing                   |
| send_context_free_inline  | Missing                   |
| send_deferred             | Deprecated, won't support |
| get_sender                | Missing                   |

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