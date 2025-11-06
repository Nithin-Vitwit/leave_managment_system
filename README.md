# 🏢 Leave Management System (Go CLI)

A simple **Command-Line Leave Management System** written in **Go**.  
It allows **Employees** to apply for leaves and view their leave history, while **HR** can review, approve, or decline pending leave applications.  

All data is stored in a local JSON file (`data.json`).

---

## 📦 Features

### 👩‍💼 Employee
- Apply for leave (with reason and date range)
- View leave history (status: Pending / Granted / Declined)
- Switch users or exit the system

### 🧑‍💼 HR
- Secure login using a password (`admin`)
- View pending leave applications
- Approve or decline leave requests
- Automatically save all changes

### 💾 System
- Automatically loads and saves data to `data.json`
- Initializes default users on first run:
  - Employees:  
    - **Nithin** (ID: `22`)  
    - **Kalyan** (ID: `23`)
  - HR:  
    - **Soumya** (ID: `12`)

---

## ⚙️ Requirements

- [Go](https://go.dev/) version **1.20+**

---

## 🚀 How to Run

1. **Clone or download** this project.
   ```bash
   git clone https://github.com/yourusername/leave-management-system.git
   cd leave-management-system
