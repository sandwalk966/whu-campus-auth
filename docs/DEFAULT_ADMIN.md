# Default Administrator Account

## 🎉 Overview

The system automatically creates a default administrator account on **first startup**, requiring no manual intervention.

---

## 👤 Default Account Information

### Administrator Account

| Field | Value |
|------|-----|
| **Username** | `admin` |
| **密码** | `your password(set in the .env file)` |
| **Nickname** | `System Administrator` |
| **Role** | `Super Administrator` |
| **Role Code** | `admin` |
| **Email** | `admin@system.local` |


## 🚀 Startup Flow

### First Startup

```
1. Start Application
   ↓
2. Initialize Database
   ↓
3. Execute Database Migration (Create Tables)
   ↓
4. Call InitAdminUser(db)
   ├─ Check if administrator role exists
   │  └─ If not → Create `admin` role
   ├─ Check if administrator account exists
   │  └─ If not → Create `admin` account
   │      └─ Assign `admin` role
   └─ Initialize dictionary data
   ↓
5. Startup Complete
```

**Log Output**:
```
[INFO] Administrator role created successfully id=1
[INFO] Default administrator account created successfully id=1 username=admin role_id=1
[INFO] Server started successfully, listening on: :8888
```

### Subsequent Startups

```
1. Start Application
   ↓
2. Initialize Database
   ↓
3. Call InitAdminUser(db)
   ├─ Check administrator role → ✅ Already exists, skip
   ├─ Check administrator account → ✅ Already exists, skip
   └─ Check role assignment → ✅ Already assigned, skip
   ↓
4. Startup Complete
```

**Log Output**:
```
[INFO] Administrator role already exists id=1
[INFO] Administrator account already exists id=1
[INFO] Server started successfully, listening on: :8888
```

---