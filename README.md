# 🏗️ 16-bit Virtual Machine Emulator

## 📌 Overview
This project is a **Virtual Machine** written in **Go** that simulates a simple instruction set architecture (ISA). It includes **memory management**, **register operations**, and **basic arithmetic instructions**.

---

## 🚀 Features
- **Memory Management** 📦
  - Supports **256x256 bytes** of memory
  - Read & write operations
- **Register Operations** 🏛️
  - Supports multiple registers (`IP`, `ACC`, `R1`, `R2`, etc.)
  - Register to register, memory to register operations
- **Instruction Set** ⚙️
  - `MOV_LIT_REG`: Move literal to register
  - `MOV_REG_REG`: Move register to register
  - `MOV_REG_MEM`: Move register value to memory
  - `MOV_MEM_REG`: Move memory value to register
  - `ADD_REG_REG`: Add two registers
- **Debugging Tools** 🔍
  - View memory at specific addresses
  - Step-by-step execution

---

## 📜 Instruction Set

| Opcode (Hex) | Instruction     | Description                            |
|-------------|----------------|----------------------------------------|
| `0x10`      | `MOV_LIT_REG`   | Move a **literal** into a **register** |
| `0x11`      | `MOV_REG_REG`   | Move a **register value** to another  |
| `0x12`      | `MOV_REG_MEM`   | Move **register value** to **memory** |
| `0x13`      | `MOV_MEM_REG`   | Move **memory value** to **register** |
| `0x14`      | `ADD_REG_REG`   | **Add** two register values           |

---

## 🏗️ Project Structure
```
📂 project-root
 ├── 📂 cpu/           # CPU implementation
 ├── 📂 memory/        # Memory management
 ├── 📂 constants/     # Instruction set definitions
 ├── 📜 main.go        # Entry point of the program
 ├── 📜 README.md      # This documentation file
```

---

## 🔧 How to Use

### 1️⃣ **Install Go**
Make sure you have Go installed. If not, download it from [golang.org](https://golang.org/dl/).

### 2️⃣ **Clone the Repository**
```sh
git clone https://github.com/Martbul/16BitVirtualMachine.git
cd 16BitVirtualMachine
```

### 3️⃣ **Run the Emulator**
```sh
go run main.go
```

---

## 🛠️ Debugging & Testing

### **View Instruction Memory Before Execution**
To check what instructions are loaded into memory, add:
```go
fmt.Println("Instruction Memory:", writableBytes[:i])
```

### **Step-by-Step Execution**
To execute and debug step by step:
```go
cpu.Step()
cpu.Debug()
cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
cpu.ViewMemoryAt(0x0100)
```

---

## 🎯 Example Execution Output
```
Instruction Memory: [16 18 52 2 16 171 205 3 20 2 3 18 1 1 0]
🔹 Step 1: MOV_LIT_REG 0x1234 → R1
🔹 Step 2: MOV_LIT_REG 0xABCD → R2
🔹 Step 3: ADD_REG_REG R1 + R2 → ACC
🔹 Step 4: MOV_REG_MEM ACC → Memory[0x0100]
```

---

## 🤝 Contributing
✅ Feel free to fork and submit pull requests!

---

## 📝 License
This project is licensed under the **MIT License**.

