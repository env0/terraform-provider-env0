# For project variables: include a resource block that passes the project_id argument
terraform import env0_configuration_variable.my_project_config '{ "Scope": "PROJECT", "ScopeId": "project id", "name": "configuration variable name"}'
# For global variables: do not pass ScopeId in the import commant
terraform import env0_configuration_variable.my_global_config '{ "Scope": "GLOBAL", "name": "configuration variable name"}'
