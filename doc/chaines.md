| Table        | PREROUTING | INPUT | FORWARD | OUTPUT | POSTROUTING |
| ------------ | ---------- | ----- | ------- | ------ | ----------- |
| **raw**      | ✅          | ❌     | ❌       | ✅      | ❌           |
| **mangle**   | ✅          | ✅     | ✅       | ✅      | ✅           |
| **nat**      | ✅          | ❌     | ❌       | ✅      | ✅           |
| **filter**   | ❌          | ✅     | ✅       | ✅      | ❌           |
| **security** | ❌          | ✅     | ❌       | ❌      | ❌           |
