# List of error codes

## Error listing

### 1: Database error.
 - [pkg/api/middleware/recovery.go:52](../pkg/api/middleware/recovery.go#L52)

### 2: Database inconsistency.
 - [pkg/api/middleware/uses_auth.go:54](../pkg/api/middleware/uses_auth.go#L54)

### 3: Invalid Authorization header.
 - [pkg/api/middleware/requires_auth.go:14](../pkg/api/middleware/requires_auth.go#L14)
 - [pkg/api/middleware/uses_auth.go:28](../pkg/api/middleware/uses_auth.go#L28)
 - [pkg/api/middleware/uses_auth.go:34](../pkg/api/middleware/uses_auth.go#L34)

### 4: Invalid JSON input.
 - [pkg/api/accounts/create.go:25](../pkg/api/accounts/create.go#L25)
 - [pkg/api/applications/create.go:28](../pkg/api/applications/create.go#L28)
 - [pkg/api/applications/update.go:50](../pkg/api/applications/update.go#L50)
 - [pkg/api/oauth/handler.go:24](../pkg/api/oauth/handler.go#L24)
 - [pkg/api/resources/create.go:24](../pkg/api/resources/create.go#L24)
 - [pkg/api/resources/update.go:49](../pkg/api/resources/update.go#L49)

### 5: Token not found.
 - [pkg/api/middleware/uses_auth.go:44](../pkg/api/middleware/uses_auth.go#L44)

### 6: Invalid token type, must be "auth".
 - [pkg/api/middleware/uses_auth.go:59](../pkg/api/middleware/uses_auth.go#L59)

### 7: Authentication token expired.
 - [pkg/api/middleware/uses_auth.go:64](../pkg/api/middleware/uses_auth.go#L64)

### 8: Account is blocked.
 - [pkg/api/middleware/uses_auth.go:69](../pkg/api/middleware/uses_auth.go#L69)

### 9: Insufficient token permissions.
 - [pkg/api/accounts/delete.go:41](../pkg/api/accounts/delete.go#L41)
 - [pkg/api/applications/create.go:39](../pkg/api/applications/create.go#L39)
 - [pkg/api/applications/delete.go:40](../pkg/api/applications/delete.go#L40)
 - [pkg/api/applications/list.go:49](../pkg/api/applications/list.go#L49)
 - [pkg/api/applications/list.go:56](../pkg/api/applications/list.go#L56)
 - [pkg/api/applications/update.go:58](../pkg/api/applications/update.go#L58)
 - [pkg/api/resources/create.go:35](../pkg/api/resources/create.go#L35)
 - [pkg/api/resources/delete.go:40](../pkg/api/resources/delete.go#L40)
 - [pkg/api/resources/update.go:57](../pkg/api/resources/update.go#L57)

### 10: Invalid callback format.
 - [pkg/api/applications/create.go:45](../pkg/api/applications/create.go#L45)
 - [pkg/api/applications/update.go:71](../pkg/api/applications/update.go#L71)

### 11: Application name is invalid.
 - [pkg/api/applications/create.go:48](../pkg/api/applications/create.go#L48)
 - [pkg/api/applications/update.go:74](../pkg/api/applications/update.go#L74)

### 12: Invalid e-mail format.
 - [pkg/api/accounts/create.go:31](../pkg/api/accounts/create.go#L31)
 - [pkg/api/applications/create.go:51](../pkg/api/applications/create.go#L51)
 - [pkg/api/applications/update.go:77](../pkg/api/applications/update.go#L77)

### 13: Invalid home page format.
 - [pkg/api/applications/create.go:54](../pkg/api/applications/create.go#L54)
 - [pkg/api/applications/update.go:80](../pkg/api/applications/update.go#L80)

### 14: Invalid ID format.
 - [pkg/api/accounts/delete.go:24](../pkg/api/accounts/delete.go#L24)
 - [pkg/api/accounts/get.go:21](../pkg/api/accounts/get.go#L21)
 - [pkg/api/accounts/get.go:31](../pkg/api/accounts/get.go#L31)
 - [pkg/api/applications/delete.go:22](../pkg/api/applications/delete.go#L22)
 - [pkg/api/applications/get.go:18](../pkg/api/applications/get.go#L18)
 - [pkg/api/applications/update.go:24](../pkg/api/applications/update.go#L24)
 - [pkg/api/resources/delete.go:22](../pkg/api/resources/delete.go#L22)
 - [pkg/api/resources/get.go:18](../pkg/api/resources/get.go#L18)
 - [pkg/api/resources/update.go:25](../pkg/api/resources/update.go#L25)

### 15: Application not found.
 - [pkg/api/applications/delete.go:32](../pkg/api/applications/delete.go#L32)
 - [pkg/api/applications/get.go:28](../pkg/api/applications/get.go#L28)
 - [pkg/api/applications/update.go:34](../pkg/api/applications/update.go#L34)
 - [pkg/api/oauth/authorization_code.go:22](../pkg/api/oauth/authorization_code.go#L22)
 - [pkg/api/oauth/password_grant.go:36](../pkg/api/oauth/password_grant.go#L36)

### 16: Invalid owner ID.
 - [pkg/api/applications/list.go:41](../pkg/api/applications/list.go#L41)

### 17: Invalid date created start.
 - [pkg/api/applications/list.go:92](../pkg/api/applications/list.go#L92)

### 18: Invalid date created end.
 - [pkg/api/applications/list.go:98](../pkg/api/applications/list.go#L98)

### 19: Invalid date modified start.
 - [pkg/api/applications/list.go:106](../pkg/api/applications/list.go#L106)

### 20: Invalid date modified end.
 - [pkg/api/applications/list.go:112](../pkg/api/applications/list.go#L112)

### 21: Invalid skip format.
 - [pkg/api/applications/list.go:200](../pkg/api/applications/list.go#L200)

### 22: Invalid limit format.
 - [pkg/api/applications/list.go:209](../pkg/api/applications/list.go#L209)

### 23: Outdated object used in PUT.
 - [pkg/api/applications/update.go:65](../pkg/api/applications/update.go#L65)
 - [pkg/api/resources/update.go:64](../pkg/api/resources/update.go#L64)

### 24: Invalid email domain.
 - [pkg/api/accounts/create.go:33](../pkg/api/accounts/create.go#L33)

### 25: Address is taken.
 - [pkg/api/accounts/create.go:39](../pkg/api/accounts/create.go#L39)

### 26: Invalid password format.
 - [pkg/api/accounts/create.go:42](../pkg/api/accounts/create.go#L42)

### 27: Account not found.
 - [pkg/api/accounts/delete.go:33](../pkg/api/accounts/delete.go#L33)
 - [pkg/api/accounts/get.go:41](../pkg/api/accounts/get.go#L41)
 - [pkg/api/oauth/password_grant.go:52](../pkg/api/oauth/password_grant.go#L52)
 - [pkg/api/oauth/password_grant.go:62](../pkg/api/oauth/password_grant.go#L62)

### 28: Invalid application secret.
 - [pkg/api/oauth/authorization_code.go:27](../pkg/api/oauth/authorization_code.go#L27)

### 29: Authorization code not found.
 - [pkg/api/oauth/authorization_code.go:37](../pkg/api/oauth/authorization_code.go#L37)
 - [pkg/api/oauth/authorization_code.go:42](../pkg/api/oauth/authorization_code.go#L42)

### 30: Invalid grant type.
 - [pkg/api/oauth/handler.go:47](../pkg/api/oauth/handler.go#L47)

### 31: Invalid expiry date.
 - [pkg/api/oauth/password_grant.go:42](../pkg/api/oauth/password_grant.go#L42)

### 32: Invalid password.
 - [pkg/api/oauth/password_grant.go:67](../pkg/api/oauth/password_grant.go#L67)

### 33: Resource not found.
 - [pkg/api/resources/delete.go:32](../pkg/api/resources/delete.go#L32)
 - [pkg/api/resources/get.go:28](../pkg/api/resources/get.go#L28)
 - [pkg/api/resources/update.go:35](../pkg/api/resources/update.go#L35)
