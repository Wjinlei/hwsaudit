#include "myacl.h"

#include <regex.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/acl.h>

char* getfacl(char* file)
{
        struct re_pattern_buffer regex;
        struct __acl_ext* facl;
        char* facl_text;
        char* delim = "\n";
        char* sub_string;
        char* result = malloc(4096 * sizeof(char));
        sprintf(result, "%s", "");

        // Get acl text
        facl = acl_get_file(file, ACL_TYPE_ACCESS);
        if (facl == NULL) {
                return result;
        }
        facl_text = acl_to_text(facl, NULL);

        // Init regex
        if (regcomp(&regex, ":.+:", REG_EXTENDED) != 0) {
                return result;
        }

        // Filter empty acl rule
        sub_string = strtok(facl_text, delim);
        while (sub_string != NULL) {
                int regex_result = regexec(&regex, sub_string, 0, NULL, 0);
                if (REG_NOERROR == regex_result) {
                        strncat(result, delim, strlen(delim));
                        strncat(result, sub_string, strlen(sub_string));
                }
                sub_string = strtok(NULL, delim);
        }

        acl_free(facl);
        acl_free(facl_text);
        return result;
}
