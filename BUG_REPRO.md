# Bug Reproduction

Bearer parsing accepts a token without a separated scheme, repository errors lose their sentinel identity, canceled permission reads return data, and a wildcard allow can override an explicit deny.

Send malformed bearer values and exercise the repository, canceled context, and wildcard/deny permission paths. Authentication or authorization returns the wrong result and the auth tests fail.

