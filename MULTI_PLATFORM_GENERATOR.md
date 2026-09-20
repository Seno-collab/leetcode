# Multi-platform generator proposal

Trạng thái: **chỉ là đề xuất, chưa triển khai**.

Khi muốn bắt đầu, gọi:

```text
Hãy apply MULTI_PLATFORM_GENERATOR.md
```

## Mục tiêu

Import một bài từ Codeforces, AtCoder, CSES, HackerRank, Kattis hoặc các online
judge phổ biến; tự tạo Go template, lưu sample input/output và chạy kiểm tra
AC/WA bằng một CLI duy nhất.

CLI dự kiến:

```bash
./lc import <problem-url>
./lc <problem-name>
./lc test <problem-name>
./lc listen
```

Ví dụ:

```bash
./lc import https://codeforces.com/problemset/problem/954/G
./lc cf_954_g
```

Kết quả generate:

```text
exercises/cf_954_g/
├── main.go
├── problem.json
└── tests/
    ├── sample-1.in
    ├── sample-1.out
    ├── sample-2.in
    └── sample-2.out
```

## Kiến trúc đã chọn

1. Dùng JSON format của Competitive Companion làm contract chung.
2. `./lc import <URL>` gọi `oj-api get-problem` để lấy metadata và sample.
3. `./lc listen` nhận JSON từ browser extension Competitive Companion.
4. Không viết scraper riêng cho từng website.
5. Cache metadata/sample đã tải để không request website lặp lại.

Model tối thiểu:

```text
Problem
├── name
├── group
├── url
├── interactive
├── timeLimit
├── memoryLimit
└── tests[]
    ├── input
    └── output
```

## Thay đổi runner bắt buộc

Runner hiện tại dùng `solve([]int) int`; cần chuyển template competitive
programming sang raw stdin/stdout:

```go
func solve(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Parse input theo format của đề.
}
```

Mỗi sample phải giữ nguyên toàn bộ nội dung nhiều dòng trong file `.in` và
`.out`. Runner build binary một lần, truyền từng `.in` qua stdin, thu stdout và
so sánh với `.out`.

## Phạm vi theo nền tảng

- Codeforces, AtCoder: ưu tiên đầu tiên qua `oj-api`.
- HackerRank, Kattis, CodeChef: dùng cùng adapter nếu `oj-api` lấy được sample.
- CSES và các judge khác: nhận qua Competitive Companion.
- LeetCode: giữ adapter riêng vì dùng function signature thay vì stdin/stdout.

## MVP

- [ ] Thêm model `Problem` và `Sample`.
- [ ] Thêm template stdin/stdout cho Go.
- [ ] Lưu sample tại `tests/sample-N.in` và `tests/sample-N.out`.
- [ ] Chạy sample và báo AC/WA từng case.
- [ ] Giữ đo thời gian xử lý từng sample và tổng thời gian.
- [ ] Thêm `./lc import <URL>` qua `oj-api`.
- [ ] Hỗ trợ Codeforces và AtCoder trước.
- [ ] Sinh slug ổn định như `cf_954_g`, `at_abc461_a`.
- [ ] Lưu URL/time limit/memory limit trong `problem.json`.
- [ ] Báo lỗi rõ ràng nếu thiếu `oj-api`.
- [ ] Không ghi đè exercise đã tồn tại.

## Sau MVP

- [ ] `./lc listen` cho Competitive Companion.
- [ ] Import cả contest.
- [ ] Custom test do người dùng thêm.
- [ ] Timeout khi chạy sample.
- [ ] Adapter LeetCode dạng function.
- [ ] Special checker khi một bài có nhiều output hợp lệ.

## Chưa hỗ trợ trong MVP

- Interactive problem.
- Special judge phức tạp.
- Tự động login hoặc submit bài.
- Tải hidden/system test không công khai.
- Tự sinh lời giải cho bài toán.

## Điều kiện hoàn thành

Một vertical slice được xem là hoàn thành khi:

1. Import được một URL Codeforces và một URL AtCoder thật.
2. Tạo đúng `main.go`, `problem.json` và các cặp sample `.in/.out`.
3. `./lc <name>` build một lần và báo AC/WA đúng cho từng sample.
4. Input nhiều dòng, chuỗi và số đều được giữ nguyên.
5. Test lỗi URL, thiếu dependency, exercise trùng tên và interactive problem.
6. `go test ./...`, shell syntax check và `git diff --check` đều pass.

## Nguồn tham khảo

- Competitive Companion format và danh sách website:
  https://github.com/jmerle/competitive-companion
- `oj-api` và danh sách online judge được hỗ trợ:
  https://github.com/online-judge-tools/api-client
- `oj` download/test sample:
  https://github.com/online-judge-tools/oj
- Codeforces API:
  https://codeforces.com/apiHelp
