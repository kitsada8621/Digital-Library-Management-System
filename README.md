# Digital Library Management System (DLMS)
DLMS: A comprehensive digital library management system that enables efficient management of digital resources, from storage and categorization to search and lending services for users.

### Installation
- clone project ไปยังเครื่องของคุณ `git clone https://github.com/kitsada8621/Digital-Library-Management-System.git`
- run คำสั่ง `go get .` เพื่อทำการติดตั้ง dependenc
- ตรวจสอบ config ต่างในไฟล์ .env ให้ถูกต้อง
- เริ่มต้นการใช้งาน project ด้วยคำสั่ง `go run main.go`

### Postman
 - link: https://drive.google.com/drive/folders/1HNNGRwmjNcm0vPMml6sNvUm3RV_wAeW2?usp=sharing

### Test
ภายหลังการเริ่มต้นระบบและการเชื่อมต่อฐานข้อมูลสำเร็จ ระบบจะจัดเตรียมบัญชีผู้ใช้เริ่มต้นจำนวนสามบัญชีที่มีบทบาทแตกต่างกัน เพื่ออำนวยความสะดวกในการทดสอบระบบ ดังนี้
| Role | Username | Password |
| ------ | ------ | ------ | 
| User | user | `Dev123456!` |
| Author | author | `Dev123456!` |
| Admin | admin | `Dev123456!` |