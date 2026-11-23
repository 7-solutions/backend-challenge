# Lottery Search & Atomic Ticket Allocation System Using MongoDB

**Version:** 1.0  
**Architecture:** MongoDB-Based Distributed Search System  
**Dataset:** 1,000,000 Lottery Tickets  
**Design Approach:** MongoDB Indexing + Bitmap Index + Atomic Reservation  
**Scope:** System Design Only (No Code Implementation)

---

## 1. Executive Summary

ระบบ Lottery Search System ถูกออกแบบมาเพื่อรองรับการค้นหาหมายเลขสลากกินแบ่งจำนวน 1 ล้านรายการ โดยแต่ละหมายเลขประกอบด้วยตัวเลข 6 หลัก และสามารถค้นหาผ่าน Search Pattern ที่มี Wildcard (`*`) ได้ทุกตำแหน่ง

ระบบต้องสามารถค้นหาและจัดสรรสลากให้ผู้ใช้งานจำนวนมากพร้อมกัน โดยป้องกันไม่ให้สลากใบเดียวกันถูกจองโดยผู้ใช้มากกว่าหนึ่งคนในช่วงเวลาเดียวกัน

**แนวทางที่เสนอ**

1. ใช้ **MongoDB Replica Set** เป็น Primary Database เพื่อเก็บข้อมูลสลากและสถานะการจอง
2. ใช้ **MongoDB Indexing** สำหรับค้นหาหมายเลขสลากด้วย Prefix, Suffix และตำแหน่งตัวเลข
3. ใช้ **In-Memory Bitmap Index** เพื่อเพิ่มประสิทธิภาพการค้นหา Wildcard ในทุกตำแหน่ง
4. ใช้ **Atomic Conditional Update (Compare-and-Set)** ของ MongoDB เพื่อป้องกันการจองซ้ำ
5. ใช้ **Reservation Lease** เพื่อกำหนดอายุการจองและคืนสลากเข้าสู่ระบบเมื่อหมดอายุ
6. ใช้ **MongoDB Change Streams** เพื่อปรับปรุงข้อมูลใน Search Cache ให้สอดคล้องกับฐานข้อมูล

หลักการสำคัญคือ **Bitmap ใช้เร่งการค้นหา แต่ MongoDB เป็นผู้ตัดสินสิทธิ์การจองเสมอ (Single Source of Truth)**

## 2. Requirements & Assumptions

### 2.1 Functional Requirements

| Requirement | Description |
|---|---|
| Ticket Storage | รองรับสลากอย่างน้อย 1,000,000 รายการ |
| Number Format | หมายเลข 6 หลัก รวมเลขศูนย์นำหน้า |
| Wildcard Search | รองรับตัวเลข 0–9 และ `*` ทุกตำแหน่ง |
| Concurrent Allocation | ผู้ใช้หลายคนค้นหาและจองพร้อมกันได้ |
| Duplicate Prevention | สลากเดียวกันต้องไม่มี Active Reservation ซ้อนกัน |
| Reservation Expiration | สลากที่หมดเวลาจองสามารถกลับมาจองได้ |
| Confirmation | สลากที่ยืนยันการซื้อแล้วไม่สามารถจองใหม่ได้ |

### 2.2 Assumptions

- หมายเลขอยู่ระหว่าง `000000` ถึง `999999`
- สลากหนึ่งรายการมี Ticket ID ที่ไม่ซ้ำกัน
- กรณีมีสลากหลายใบที่ใช้หมายเลขเดียวกัน ให้แยกด้วย Ticket ID
- การจองมีอายุเริ่มต้น 5 นาที (สามารถปรับตาม Business Requirements)
- ระบบรองรับการค้นหาและการจองจาก Application Server หลาย Instance
- การจองต้องตรวจสอบสถานะจริงจาก MongoDB ก่อนยืนยันทุกครั้ง

**ข้อสังเกต:** หมายเลข 6 หลักมีความเป็นไปได้ทั้งหมด $10^6 = 1,000,000$ หมายเลข หากกำหนดให้หมายเลขไม่ซ้ำกันและมีข้อมูลครบ 1 ล้านรายการ เท่ากับมีหมายเลขครบทุกความเป็นไปได้

## 3. Proposed System Architecture

### 3.1 High-Level Architecture

ระบบแบ่งออกเป็น 5 ส่วนหลัก

| Component | Responsibility |
|---|---|
| API Gateway / Load Balancer | รับคำขอและกระจายไปยัง Application Instance |
| Lottery Search Service | ตรวจสอบ Pattern และคำนวณ Candidate Tickets |
| Bitmap Index / Search Cache | ค้นหา Ticket ID ที่ตรง Pattern อย่างรวดเร็ว |
| Allocation Service | จองสลากด้วย Atomic Operation |
| MongoDB Replica Set | เก็บข้อมูลสถานะจริงและรับประกันความถูกต้องของการจอง |

**Application Flow**

Client → API Gateway → Lottery Search Service → Bitmap Search → Candidate Ticket IDs → Allocation Service → MongoDB Atomic Update → Reservation Result

**Background Processes**

- MongoDB Change Streams → Search Cache Synchronization
- Reservation Expiration Worker → Release Expired Reservations
- Monitoring → Metrics, Logs และ Alerts

### 3.2 Architectural Decision

เลือกใช้ MongoDB เป็นฐานข้อมูลหลัก เนื่องจากรองรับ Document-Based Data Model, Secondary Indexes, Atomic Single-Document Updates และ Replication

การทำ Atomic Update ใน Document เดียวเหมาะกับการจองสลาก เพราะสถานะการจองและข้อมูลเจ้าของการจองสามารถเก็บไว้ใน Ticket Document เดียวกัน ทำให้ไม่จำเป็นต้องใช้ Distributed Lock หรือ Multi-Document Transaction สำหรับการป้องกันสลากซ้ำในกรณีพื้นฐาน

MongoDB Replica Set จำนวน 3 Data-Bearing Nodes เป็นจุดเริ่มต้นที่เหมาะสม โดยสามารถขยายจำนวน Application Instances ได้ตามภาระงาน

## 4. Data Model Design

### 4.1 Collection: `tickets`

จัดเก็บหมายเลขสลากและสถานะการจัดสรร

| Field | Data Type | Description |
|---|---|---|
| `_id` | ObjectId / Integer | Unique Ticket Identifier |
| `ticketNo` | String | หมายเลขสลาก 6 หลัก เช่น `001235` |
| `reversedNo` | String | หมายเลขย้อนกลับสำหรับ Suffix Search |
| `d0`–`d5` | Integer | ตัวเลขแต่ละตำแหน่ง |
| `status` | String | AVAILABLE, RESERVED, SOLD |
| `reservedBy` | String / Null | User ID ของผู้จอง |
| `reservationId` | String / Null | Unique Reservation Token |
| `reservedUntil` | Date / Null | เวลาหมดอายุการจอง |
| `version` | Integer | Version สำหรับ Optimistic Concurrency |
| `createdAt` | Date | วันที่สร้างข้อมูล |
| `updatedAt` | Date | วันที่ปรับปรุงล่าสุด |

ตัวอย่างข้อมูลเชิงแนวคิด

| Field | Value |
|---|---|
| `_id` | 125 |
| `ticketNo` | `001235` |
| `reversedNo` | `532100` |
| `d0`–`d5` | 0, 0, 1, 2, 3, 5 |
| `status` | AVAILABLE |
| `reservedBy` | Null |
| `reservedUntil` | Null |

การเก็บหมายเลขเป็น String ช่วยรักษาเลขศูนย์นำหน้า ส่วนการแยกหลักตัวเลขช่วยสนับสนุนการค้นหาแบบไม่ต่อเนื่อง

หากหนึ่งหมายเลขมีสลากหลายใบ ห้ามกำหนด Unique Index บน `ticketNo` เพียงอย่างเดียว แต่ให้ใช้ `_id` หรือ Ticket ID เป็นตัวระบุสลากที่ไม่ซ้ำกัน

### 4.2 Collection: `allocation_requests`

ใช้จัดเก็บคำขอจัดสรรสลากเพื่อรองรับการ Retry และ Idempotency

| Field | Description |
|---|---|
| `requestId` | Unique Request Identifier |
| `userId` | ผู้ร้องขอ |
| `pattern` | Pattern ที่ค้นหา |
| `ticketId` | สลากที่ได้รับการจัดสรร |
| `reservationId` | Reservation Token |
| `status` | PENDING, COMPLETED, FAILED |
| `createdAt` | เวลาสร้างคำขอ |

กำหนด Unique Index บน `requestId` เพื่อให้ Client Retry ด้วย Request ID เดิมแล้วสามารถรับผลลัพธ์เดิมได้ โดยต้องมีกลไก Recovery สำหรับกรณีการจองสำเร็จแต่การบันทึก Request Result ล้มเหลว

หากต้องการรับประกันผลลัพธ์แบบ Idempotent อย่างเข้มงวดข้ามสอง Collection สามารถใช้ MongoDB Transaction เพื่อบันทึกผลการจองและ Request Result ร่วมกัน

## 5. Wildcard Search Algorithm

### 5.1 Pattern Matching

ระบบรองรับ Pattern ที่มีความยาว 6 ตัวอักษร โดยแต่ละตำแหน่งเป็นตัวเลข 0–9 หรือ `*`

| Pattern | Matching Rule | ตัวอย่างผลลัพธ์ |
|---|---|---|
| `****23` | ลงท้ายด้วย 23 | 000023, 145623, 999923 |
| `1****5` | ขึ้นต้นด้วย 1 และลงท้ายด้วย 5 | 100005, 123455, 199995 |
| `123***` | ขึ้นต้นด้วย 123 | 123000, 123456, 123999 |
| `1*3*5*` | ตำแหน่ง 1, 3, 5 ตรงตามเงื่อนไข | 103050, 123456 |
| `******` | ตรงกับทุกหมายเลข | 000000–999999 |
| `123456` | หมายเลขตรงทุกหลัก | 123456 |

### 5.2 Search Strategy A: MongoDB Native Indexing

เป็นแนวทางพื้นฐานที่สามารถใช้งานได้โดยไม่ต้องพึ่ง Search Engine ภายนอก

**Prefix Search** — ตัวอย่าง `123***`

ใช้ข้อมูล `ticketNo` ร่วมกับ B-Tree Index เพื่อค้นหาหมายเลขที่ขึ้นต้นด้วย `123`

Anchored Prefix Regex สามารถใช้ประโยชน์จาก Index Range ได้ จึงเหมาะกับ Pattern ที่กำหนดเลขจากด้านซ้ายต่อเนื่องกัน

**Suffix Search** — ตัวอย่าง `****23`

ใช้ Field `reversedNo` ซึ่งเก็บหมายเลขย้อนกลับ และค้นหาด้วย Prefix `32` บน Field ดังกล่าว

**Arbitrary Wildcard Search** — ตัวอย่าง `1*3*5*`

ใช้ Field `d0`, `d2`, `d4` เพื่อค้นหาหมายเลขที่มีตัวเลขตามตำแหน่งที่กำหนด

สร้าง Index บน Field ที่ใช้ค้นหาบ่อย โดยอาจใช้ Composite Index สำหรับ Pattern ที่พบเป็นประจำ

ข้อจำกัดคือ MongoDB ไม่สามารถใช้ Compound Index เดียวเพื่อค้นหาทุก Combination ของตำแหน่งได้อย่างมีประสิทธิภาพเสมอไป โดยเฉพาะเมื่อมีช่องว่างระหว่างตำแหน่งที่กำหนด

ดังนั้นจึงเสนอ Bitmap Index เป็น Search Accelerator สำหรับ Pattern ที่ซับซ้อน

### 5.3 Search Strategy B: Bitmap Inverted Index (Recommended Optimization)

ใช้ Bitmap เพื่อเก็บความสัมพันธ์ระหว่างตำแหน่งตัวเลขกับ Ticket ID

สร้าง Bitmap ตามตำแหน่งและค่าตัวเลข เช่น

- Position 0, Digit 1
- Position 1, Digit 3
- Position 2, Digit 5
- Position 5, Digit 9

แต่ละ Bitmap มี Bit สำหรับ Ticket ID แต่ละรายการ โดย Bit = 1 หมายถึงสลากตรงกับเงื่อนไขของ Bitmap นั้น

**จำนวน Bitmap**

- 6 ตำแหน่ง
- 10 ค่าตัวเลขต่อหนึ่งตำแหน่ง
- รวม 60 Bitmaps

สำหรับ Dataset จำนวน 1 ล้านรายการ Bitmap แบบ Dense ใช้หน่วยความจำพื้นฐานประมาณ 7.5 MB สำหรับ 60 Bitmaps โดยไม่รวม Overhead ของโครงสร้างข้อมูลและข้อมูลสถานะ

สามารถใช้ Roaring Bitmap หรือ Compressed Bitmap เพื่อลดพื้นที่หน่วยความจำและเพิ่มประสิทธิภาพการทำ Set Operations ตามลักษณะการกระจายข้อมูลจริง

**ขั้นตอนการค้นหา** — ตัวอย่าง Pattern `1*3*5*`

1. อ่าน Bitmap ของ Position 0 = 1
2. อ่าน Bitmap ของ Position 2 = 3
3. อ่าน Bitmap ของ Position 4 = 5
4. ทำ Bitwise AND ของทั้งสาม Bitmap
5. นำผลลัพธ์มา Intersect กับ Availability Bitmap
6. เลือก Candidate Ticket ID
7. ตรวจสอบและจองผ่าน MongoDB Atomic Update

การใช้ AND จะให้เฉพาะ Ticket ID ที่ตรงกับทุกเงื่อนไข โดยไม่ต้องอ่าน Ticket Document ทั้ง 1 ล้านรายการมาทำ Pattern Matching ทีละรายการ

**ข้อควรระวัง:** Availability Bitmap เป็นข้อมูลช่วยค้นหา ไม่ใช่แหล่งข้อมูลสถานะที่มีผลผูกพันทางธุรกิจ หาก Cache ไม่ทันสมัยหรือไม่พบ Candidate ต้องมี MongoDB Fallback และกลไก Cache Reconciliation เพื่อหลีกเลี่ยงการรายงานว่าสลากหมดทั้งที่ยังมีสลากว่าง

### 5.4 Search Selectivity

หากมีหมายเลขครบทุกค่าตั้งแต่ `000000` ถึง `999999` และไม่มีหมายเลขซ้ำ สามารถคำนวณจำนวนผลลัพธ์ก่อนหักสลากที่ถูกจองหรือขายแล้วได้ดังนี้

$$
Matches = 10^{6-k}
$$

โดย $k$ คือจำนวนตำแหน่งที่ระบุตัวเลขแน่นอน

| Pattern | Fixed Digits | Matching Tickets |
|---|---:|---:|
| `******` | 0 | 1,000,000 |
| `1*****` | 1 | 100,000 |
| `****23` | 2 | 10,000 |
| `1****5` | 2 | 10,000 |
| `123***` | 3 | 1,000 |
| `1*3*5*` | 3 | 1,000 |
| `1234**` | 4 | 100 |
| `12345*` | 5 | 10 |
| `123456` | 6 | 1 |

สำหรับข้อมูลที่ไม่ครอบคลุมทุกหมายเลข หรือมีหมายเลขซ้ำ จำนวนที่ตรงจริงต้องคำนวณจาก Dataset

## 6. MongoDB Indexing Strategy

กำหนด Index ตาม Query Pattern ที่เกิดขึ้นจริง แทนการสร้าง Compound Index ครอบคลุมทุก Combination

| Index | Purpose |
|---|---|
| `_id` Unique | Direct Ticket Lookup และ Atomic Allocation |
| `ticketNo` | Exact / Prefix Search |
| `reversedNo` | Suffix Search |
| `(status, d0)` | ค้นหาตามหลักแรกและสถานะ |
| `(status, d5)` | ค้นหาตามหลักสุดท้ายและสถานะ |
| `(status, d0, d5)` | Optimize Pattern เช่น `1****5` |
| `(status, reservedUntil)` | ค้นหาการจองที่หมดอายุ |
| `allocation_requests.requestId` Unique | Prevent Duplicate Request Processing |

Index สำหรับตำแหน่ง `d1`–`d4` สามารถเพิ่มได้ตาม Query Frequency และผลจาก Query Profiling

**Index Design Principles**

- ใช้ Index ที่มี Cardinality และ Selectivity เหมาะสม
- ใช้ Compound Index สำหรับ Pattern ที่พบเป็นประจำ
- ใช้ Partial Index เมื่อเหมาะสมเพื่อลดขนาด Index
- หลีกเลี่ยงการสร้าง Index มากเกินไป เพราะทุกการ Update สถานะอาจมีต้นทุนด้าน Index Maintenance
- ตรวจสอบ Query Plan ผ่าน `explain` และ Metrics จริงก่อนปรับ Index

การสร้าง Index สำหรับทุก Pattern Combination ไม่ใช่ทางเลือกแรก เนื่องจากมี Wildcard Mask ได้ 64 รูปแบบต่อหมายเลข หาก Materialize ทุก Mask จะเกิดประมาณ 64 ล้าน Pattern-to-Ticket Associations สำหรับข้อมูล 1 ล้านรายการ ส่งผลต่อ Storage และ Write Amplification

## 7. Concurrent Allocation & Duplicate Prevention

### 7.1 Problem: Race Condition

สมมติผู้ใช้งาน A และ B ค้นหาด้วย Pattern `****23` พร้อมกัน

ทั้งสองอาจพบสลากหมายเลข `123423` จากผลการค้นหาเดียวกัน

หากระบบแยกการอ่านสถานะออกจากการบันทึกการจอง โดยไม่ตรวจสอบเงื่อนไขแบบ Atomic ผู้ใช้ทั้งสองอาจได้รับสลากเดียวกัน

**Solution: Atomic Compare-and-Set**

### 7.2 Allocation Algorithm

ขั้นตอนการจัดสรร

1. Search Service สร้าง Candidate Ticket IDs จาก Pattern
2. Allocation Service เลือก Candidate หนึ่งรายการ
3. ส่ง Conditional Atomic Update ไปยัง MongoDB โดยระบุ Ticket ID และเงื่อนไขสถานะ
4. MongoDB อนุญาตให้ Update เฉพาะเมื่อ Ticket ยัง AVAILABLE หรือ Reservation เดิมหมดอายุแล้ว
5. เมื่อสำเร็จ เปลี่ยนสถานะเป็น RESERVED พร้อมบันทึก User ID, Reservation Token และ Expiration Time
6. หาก Update ไม่สำเร็จ แสดงว่ามีผู้ใช้อื่นได้สิทธิ์หรือ Ticket ไม่สามารถจองได้ ให้เลือก Candidate ใหม่
7. ส่งผลลัพธ์กลับหลังได้รับการยืนยันการบันทึกจาก MongoDB

สามารถใช้ `findOneAndUpdate` ซึ่งรองรับการ Update Document แบบ Atomic โดยให้เงื่อนไขการตรวจสอบและการเปลี่ยนสถานะอยู่ในการดำเนินการเดียวกัน

**Invariant สำคัญ**

> สลากหนึ่ง Ticket ID ต้องมี Active Reservation ที่ถูกต้องได้ไม่เกินหนึ่งรายการในเวลาเดียวกัน

MongoDB Single-Document Atomicity ทำให้ผู้ใช้งานหลายคนไม่สามารถเปลี่ยนสถานะจาก AVAILABLE เป็น RESERVED สำหรับ Ticket เดียวกันได้สำเร็จพร้อมกัน

### 7.3 Concurrent Allocation Example

| Time | User A | User B | MongoDB Ticket Status |
|---|---|---|---|
| T1 | ค้นหา `****23` | ค้นหา `****23` | AVAILABLE |
| T2 | เลือก Ticket X | เลือก Ticket X | AVAILABLE |
| T3 | Atomic Update Success | ส่ง Atomic Update | RESERVED by A |
| T4 | ได้รับ Ticket X | Atomic Update Failed | RESERVED by A |
| T5 | — | เลือก Ticket Y | Ticket Y AVAILABLE |
| T6 | — | Atomic Update Success | Ticket Y RESERVED by B |

**Result:** User A ได้ Ticket X และ User B ได้ Ticket Y ไม่มีการจอง Ticket X ซ้ำ

แม้คำขอจะใช้ Search Pattern ต่างกัน แต่ถ้าเลือก Ticket เดียวกัน MongoDB ก็ยังป้องกันการจองซ้ำได้ เพราะการตรวจสอบสิทธิ์ผูกกับ Ticket ID ไม่ใช่ Pattern

### 7.4 Why Not Distributed Lock?

ไม่แนะนำให้ใช้ Distributed Lock ระดับ Search Pattern เช่น Lock Key ของ `****23` เป็นกลไกหลัก เพราะจะทำให้คำขอที่ใช้ Pattern เดียวกันต้องรอคิวโดยไม่จำเป็น

นอกจากนี้ Pattern ต่างกันสามารถพบสลากหมายเลขเดียวกัน เช่น `****23` และ `1****3` ซึ่ง Pattern-Level Lock เพียงอย่างเดียวไม่สามารถป้องกันการจองสลากซ้ำข้าม Pattern ได้

การทำ Atomic Allocation ที่ Ticket Document จึงเหมาะสมกว่า เพราะล็อกในระดับทรัพยากรที่ต้องการจัดสรรจริง และเปิดโอกาสให้จองสลากคนละใบพร้อมกันได้

## 8. Reservation Lifecycle

กำหนด Ticket Status ดังนี้

| Current State | Event | Next State |
|---|---|---|
| AVAILABLE | Reservation Successful | RESERVED |
| RESERVED | Confirm Purchase | SOLD |
| RESERVED | Reservation Expired | AVAILABLE |
| RESERVED | User Cancelled | AVAILABLE |
| SOLD | Normal Operation | SOLD |

### 8.1 Reservation Expiration

ทุกการจองจะมี `reservedUntil` กำหนดเวลาหมดอายุ

ระบบสามารถใช้ Background Worker เพื่อตรวจสอบการจองที่หมดอายุ และเปลี่ยนกลับเป็น AVAILABLE

เพื่อป้องกัน Race Condition ระหว่าง Worker กับผู้จองรายใหม่ การ Release ต้องเป็น Conditional Atomic Update ซึ่งตรวจสอบว่าสถานะยังเป็น RESERVED และเวลาจองยังหมดอายุจริง ณ ขณะ Update

การจองใหม่ต้องตรวจสอบเงื่อนไขหมดอายุได้ด้วย แม้ Worker ยังไม่ทำงาน เพื่อไม่ให้ Ticket ถูก Block โดยไม่จำเป็น

ไม่ควรใช้ MongoDB TTL Index ลบข้อมูลจาก `tickets` โดยตรง เนื่องจาก TTL Index ใช้สำหรับลบ Document ที่หมดอายุ ไม่ใช่เปลี่ยนสถานะ Document จาก RESERVED เป็น AVAILABLE

### 8.2 Confirmation Safety

เมื่อผู้ใช้ยืนยันการซื้อ ระบบต้องตรวจสอบพร้อมกันว่า

- Ticket ID ตรงกับที่จอง
- Status ยังเป็น RESERVED
- Reservation Token ตรงกับของผู้ใช้งาน
- Reservation ยังไม่หมดอายุ

หาก Reservation หมดอายุและถูกจัดสรรให้ผู้ใช้อื่นแล้ว ผู้จองเดิมต้องไม่สามารถยืนยันรายการได้

การเปลี่ยน RESERVED → SOLD ต้องเป็น Atomic Conditional Update เช่นเดียวกับขั้นตอนการจอง

### 8.3 Failure Recovery

รองรับสถานการณ์ Network Timeout, Application Crash และ Database Failover

ใช้ Request ID / Idempotency Key เพื่อระบุคำขอเดิม หากคำขอสำเร็จใน MongoDB แต่ Response สูญหาย ระบบต้องตรวจสอบผลเดิมก่อนทำการจัดสรรครั้งใหม่

สำหรับ Production ให้ใช้ `writeConcern: majority` เพื่อเพิ่มความทนทานต่อการเปลี่ยน Primary และใช้เวลาอ้างอิงที่เชื่อถือได้สำหรับตรวจสอบ Lease Expiration

## 9. Search Cache Consistency

Bitmap Index สามารถอยู่ใน Memory ของ Application Instance เพื่อให้ Search Service ใช้ค้นหาได้โดยตรง

**Synchronization Strategy**

1. โหลดข้อมูล Bitmap จาก MongoDB เมื่อเริ่มต้นระบบ
2. Subscribe MongoDB Change Streams เพื่อรับการเปลี่ยนแปลงของ Ticket
3. เมื่อมีการจองสำเร็จ ให้ปรับปรุง Availability Bitmap
4. เมื่อหมดอายุหรือยกเลิกการจอง ให้เพิ่ม Ticket กลับเข้าสู่ Available Candidates
5. หาก Change Streams ขาดช่วง ให้ Resume จาก Token หรือ Rebuild Cache เมื่อไม่สามารถ Resume ได้
6. หาก Cache ไม่พบสลากหรือเกิดความขัดแย้งซ้ำจำนวนมาก ให้ Fallback ไปยัง MongoDB

ระบบสามารถมี Bitmap แยกกันในแต่ละ Application Instance ได้ เนื่องจาก Bitmap ไม่ใช่กลไกบังคับสิทธิ์การจอง

**Correctness Guarantee**

- Cache อาจมีข้อมูลล่าช้าได้
- Cache ไม่สามารถยืนยันการจองได้เอง
- MongoDB Atomic Update เป็นจุดตัดสินว่าการจองสำเร็จหรือไม่
- การอ่านข้อมูล Cache ที่ไม่ทันสมัยต้องไม่ทำให้เกิด Duplicate Reservation
- ต้องมี Fallback และ Reconciliation เพื่อป้องกัน False-Negative Search Result

## 10. Performance Analysis

### 10.1 Search Complexity

กำหนดให้

- $N$ = จำนวนสลากทั้งหมด
- $k$ = จำนวนตำแหน่งที่ระบุตัวเลขแน่นอน
- $w$ = จำนวน Bits ที่ประมวลผลได้ต่อหนึ่ง Machine Word
- $C$ = จำนวน Candidate ที่ต้องตรวจสอบ

| Approach | Approximate Complexity | Tradeoff |
|---|---|---|
| Full Collection Scan | $O(6N)$ | Implementation ง่าย แต่ต้นทุนเพิ่มตามจำนวนข้อมูล |
| MongoDB Indexed Search | $O(\log N + C)$ โดยประมาณ | ขึ้นอยู่กับ Index, Query Plan และ Selectivity |
| Bitmap Intersection | $O(kN/w)$ สำหรับ Dense Bitmap แบบ Word-wise | ใช้ RAM เพิ่ม แต่ค้นหา Pattern แบบกระจายได้ |
| Atomic Allocation by Ticket ID | $O(\log N)$ โดยประมาณสำหรับ Indexed Lookup | มี Write/Replication Cost และ Retry เมื่อเกิด Contention |

Bitmap Intersection ไม่จำเป็นต้องดึง Document ทั้งหมดออกจาก MongoDB จึงช่วยลด Database Read Load ได้อย่างมากในกรณี Wildcard Search ที่ซับซ้อน

อย่างไรก็ตาม Performance จริงขึ้นอยู่กับ Availability Ratio, จำนวน Retry, Concurrent Requests, Cache Synchronization และ Hardware ของระบบ

### 10.2 Memory & Storage

สำหรับ Dataset 1 ล้านรายการ

| Component | Estimated Usage |
|---|---|
| 60 Dense Position Bitmaps | ~7.5 MB Raw |
| Availability Bitmap | ~0.125 MB Raw |
| MongoDB Ticket Documents | ขึ้นอยู่กับ BSON Schema และข้อมูลประกอบ |
| MongoDB Indexes | ขึ้นอยู่กับจำนวน Index และขนาด Index Keys |
| Reservation Logs | เพิ่มขึ้นตามจำนวนธุรกรรม |

Bitmap มีขนาดค่อนข้างเล็กเมื่อเทียบกับหน่วยความจำของ Application Server ในปัจจุบัน จึงสามารถใช้เป็น Search Accelerator แบบ In-Memory ได้

ตัวเลขด้านบนไม่รวม Memory Overhead, Data Structures, Buffer, Process Runtime หรือ Replication

### 10.3 Contention Optimization

เมื่อมีคำขอจำนวนมากสำหรับ Pattern เดียวกัน ไม่ควรเลือก Candidate แรกตามลำดับ Ticket ID เสมอไป เพราะจะทำให้ผู้ใช้หลายคนแข่งขันกันเพื่อจอง Ticket เดียวกัน

เสนอให้ใช้

- Randomized Candidate Selection
- Rotating Cursor หรือ Round-Robin Candidate Selection
- Local Best-Effort Candidate Exclusion
- Bounded Retry พร้อมเลือก Candidate ใหม่
- MongoDB Fallback เมื่อ Cache มีข้อมูลไม่เพียงพอ

การกระจาย Candidate ช่วยลด Write Contention แต่ไม่ใช่ตัวรับประกันความถูกต้อง กลไกป้องกันการจองซ้ำยังคงเป็น Atomic Update บน MongoDB

## 11. Production Deployment Recommendation

### 11.1 Infrastructure

| Component | Recommendation |
|---|---|
| Primary Database | MongoDB Replica Set |
| Database Nodes | 3 Data-Bearing Nodes |
| Storage Engine | WiredTiger |
| Search Cache | In-Memory Bitmap / Roaring Bitmap |
| Application | Multiple Stateless Service Instances |
| Load Balancing | API Gateway / Load Balancer |
| Cache Synchronization | MongoDB Change Streams |
| Reservation Processing | Atomic Conditional Update |
| Monitoring | Metrics, Logs, Distributed Tracing |

สำหรับ 1 ล้านรายการ แนะนำให้เริ่มจาก Replica Set โดยยังไม่ใช้ Sharding เนื่องจากขนาดข้อมูลไม่ได้บังคับให้ต้องกระจายข้อมูลข้าม Shard

ให้พิจารณา Sharding เมื่อการทดสอบพบข้อจำกัดด้าน Throughput, Storage หรือการขยายระบบจริง ไม่ใช่จากจำนวน Document เพียงอย่างเดียว

### 11.2 Availability & Reliability

ระบบควรมี

- Health Check และ Readiness Check
- Connection Pooling
- Request Timeout และ Bounded Retry
- Majority-Acknowledged Writes
- Database Backup และ Recovery Plan
- Reservation Expiration Monitoring
- Cache Rebuild Procedure
- Idempotent Request Handling

### 11.3 Monitoring Metrics

| Metric | Purpose |
|---|---|
| Search Latency P50/P95/P99 | วัดความเร็วในการค้นหา |
| Reservation Latency P95/P99 | วัดความเร็วในการจอง |
| Allocation Success Rate | วัดสัดส่วนการจองสำเร็จ |
| Atomic Update Conflict Rate | วัดการแข่งขันในการจอง |
| Average Retry Count | ประเมิน Write Contention |
| Cache Hit / Fallback Rate | ประเมินประสิทธิภาพ Bitmap |
| Change Stream Lag | ตรวจสอบความทันสมัยของ Cache |
| Expired Reservation Backlog | ตรวจสอบประสิทธิภาพ Cleanup Worker |
| Duplicate Active Reservations | ต้องเท่ากับศูนย์ |

## 12. Testing & Evaluation Plan

กำหนดการทดสอบเพื่อประเมินความถูกต้องและประสิทธิภาพก่อนใช้งาน Production

| Test Case | Scenario | Expected Result |
|---|---|---|
| TC01 | ค้นหา `****23` | ได้เฉพาะหมายเลขลงท้าย 23 |
| TC02 | ค้นหา `1****5` | ได้เฉพาะหมายเลขขึ้นต้น 1 ลงท้าย 5 |
| TC03 | ค้นหา `123***` | ได้เฉพาะหมายเลขขึ้นต้น 123 |
| TC04 | ค้นหา `******` | ค้นหาได้จากทุกหมายเลขที่ยังมีสิทธิ์จัดสรร |
| TC05 | ผู้ใช้ 100 คนจอง Pattern เดียวกันพร้อมกัน | ไม่พบ Ticket ID ซ้ำใน Active Reservations |
| TC06 | ผู้ใช้ 100 คนจอง Pattern ต่างกันแต่มีผลลัพธ์ทับซ้อน | ไม่มีการจอง Ticket ID เดียวกันพร้อมกัน |
| TC07 | Reservation หมดอายุ | Ticket กลับมาจองได้ |
| TC08 | ผู้จองเดิม Confirm หลัง Lease หมดอายุ | ต้องไม่สามารถ Confirm ได้ |
| TC09 | Application Instance ล้มเหลวหลังจองสำเร็จ | Retry ด้วย Request ID เดิมสามารถกู้คืนผลลัพธ์ได้ |
| TC10 | Cache มีข้อมูลล่าช้า | ไม่มี Duplicate Reservation และ Fallback ทำงานได้ |
| TC11 | MongoDB Primary Failover | ระบบกู้คืนการทำงานโดยไม่เกิด Double Allocation |
| TC12 | ทุก Ticket ที่ตรง Pattern ถูกจองหรือขายแล้ว | คืน No Availability อย่างถูกต้อง |

การทดสอบ Performance ควรวัดภายใต้ Concurrent Users หลายระดับ เช่น 100, 500 และ 1,000 คน พร้อมทดสอบทั้ง Availability สูงและ Availability ต่ำ

ตัวชี้วัดสำคัญ ได้แก่ Search Latency, Allocation Throughput, Retry Rate, Write Conflicts และ Database Resource Utilization

ค่าประสิทธิภาพเป้าหมายควรกำหนดจากการทำ Benchmark บน Infrastructure จริง ไม่ควรสรุปว่า Lookup จะใช้เวลาแน่นอนกี่ Milliseconds จาก Time Complexity เพียงอย่างเดียว

## 13. Alternative Solutions & Tradeoffs

| Approach | Advantages | Disadvantages | Recommendation |
|---|---|---|---|
| MongoDB + Regex Only | ง่ายและดูแลรักษาน้อย | Arbitrary Wildcards บางรูปแบบอาจมีการ Scan จำนวนมาก | เหมาะกับ Prototype |
| MongoDB + Positional Indexes | ค้นหาได้ทุกตำแหน่งและใช้ Native MongoDB | ต้องเลือก Index ให้เหมาะกับ Workload | เหมาะกับ Production ระยะแรก |
| MongoDB + Bitmap Index | ลดการอ่าน Database และค้นหา Pattern ได้รวดเร็ว | ต้องบริหาร Cache และ Synchronization | เหมาะกับ High-Concurrency Workload |
| MongoDB + Precomputed Pattern Keys | Exact Lookup ทุก Wildcard Mask | เพิ่ม Index Entries และ Write Amplification มาก | ไม่ใช่ทางเลือกหลัก |
| MongoDB + Distributed Lock | สามารถ Serialize การจัดสรรได้ | เพิ่ม Latency, Lock Contention และ Operational Complexity | ไม่จำเป็นสำหรับ Single-Ticket Allocation |

**Selected Design: MongoDB + Bitmap Search Accelerator + Atomic Reservation**

เหตุผลในการเลือกคือแยกประสิทธิภาพการค้นหาออกจากการรับประกันความถูกต้องของข้อมูล ทำให้ระบบสามารถ Scale การค้นหาได้โดยไม่ลดความปลอดภัยด้าน Concurrent Allocation

สำหรับการเริ่มต้นที่ต้องการความเรียบง่าย สามารถใช้ MongoDB Positional Indexing และ Atomic Allocation ก่อน แล้วเพิ่ม Bitmap Index เมื่อพบว่าการค้นหาเป็น Bottleneck

## 14. Conclusion

ระบบ Lottery Search System ที่ออกแบบนี้รองรับการค้นหาสลาก 6 หลักจำนวน 1 ล้านรายการ ด้วย Wildcard Pattern ได้ในทุกตำแหน่ง โดยใช้ MongoDB เป็น Primary Database

การออกแบบประกอบด้วยองค์ประกอบสำคัญ 3 ส่วน

**1. Efficient Searching**

ใช้ MongoDB Native Indexing สำหรับการค้นหาทั่วไป และใช้ Bitmap Inverted Index เพื่อเร่งการค้นหา Wildcard ที่ซับซ้อน ลดการ Scan และลดภาระงานบน Database

**2. Correct Concurrent Allocation**

ใช้ MongoDB Atomic Conditional Update เพื่อตรวจสอบและเปลี่ยนสถานะการจองใน Operation เดียว ป้องกัน Double Allocation แม้ผู้ใช้จะค้นหาพร้อมกันหรือใช้ Pattern ต่างกัน

**3. Production Reliability**

ใช้ Reservation Lease, Request Idempotency, Majority Writes, Replica Set, Change Streams และ Cache Reconciliation เพื่อรองรับสถานการณ์จริง เช่น Concurrent Requests, Network Failures, Expired Reservations และ Application Failover

**Final Recommendation**

ใช้ **MongoDB Replica Set + Positional/Compound Indexes + In-Memory Bitmap Index + Atomic Conditional Update** โดยให้ MongoDB เป็นแหล่งข้อมูลที่เชื่อถือได้เพียงแห่งเดียวสำหรับสถานะการจอง

แนวทางนี้ตอบโจทย์ทั้งความเร็วในการค้นหา ความถูกต้องภายใต้ Concurrent Requests และความสามารถในการขยายระบบ โดยไม่จำเป็นต้องพึ่ง Distributed Lock สำหรับการจองสลากหนึ่งรายการ

---

## References

1. MongoDB — [Atomicity and Transactions](https://www.mongodb.com/docs/manual/core/write-operations-atomicity/)
2. MongoDB — [Compound Indexes](https://www.mongodb.com/docs/manual/core/indexes/index-types/index-compound/)
3. MongoDB — [Regular Expression Query](https://www.mongodb.com/docs/manual/reference/operator/query/regex/)
4. MongoDB — [Change Streams](https://www.mongodb.com/docs/manual/changestreams/)
5. MongoDB — [TTL Indexes](https://www.mongodb.com/docs/manual/core/index-ttl/)
6. MongoDB — [Write Concern](https://www.mongodb.com/docs/manual/reference/write-concern/)
