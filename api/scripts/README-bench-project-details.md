# Benchmark: `/sentinel/projects/:id/details`

ใช้สคริปต์ `api/scripts/bench_project_details.sh` เพื่อเทียบเวลา response ก่อน/หลัง optimization

## เตรียม

1. ให้ API รันอยู่ (เช่น `make up-api`)
2. เตรียม JWT token ที่เรียก endpoint ได้
3. เลือก project id หรือ project code ที่มีข้อมูลจริง

## รัน

```bash
chmod +x api/scripts/bench_project_details.sh

./api/scripts/bench_project_details.sh \
  --base-url http://localhost:8080/api/v1 \
  --project <project-id-or-code> \
  --token "<jwt-token>" \
  --runs 20 \
  --tasks-limit 600
```

## เปรียบเทียบก่อน/หลัง

- รันครั้งที่ 1 บนโค้ด/DB ก่อนเปลี่ยน
- apply migration + deploy โค้ดใหม่
- รันครั้งที่ 2 ด้วย args เดิมให้เหมือนเดิม
- เทียบค่า `p50`, `p95`, `avg`

## หมายเหตุ

- ถ้าจะทดสอบ profile project ใหญ่ ให้เพิ่ม `--tasks-limit`
- ควรรันบน environment เดียวกันและเวลาที่ load ใกล้เคียงกัน
