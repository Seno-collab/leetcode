# go-go


 Câu 1 (buffered vs unbuffered channel) bạn có thể trả lời ngắn gọn như này:

  1. Unbuffered channel (make(chan T))

  - Không có hàng đợi.
  - send sẽ block cho đến khi có receive tương ứng.
  - Dùng khi cần đồng bộ hóa chặt giữa goroutines (handoff trực tiếp).

  2. Buffered channel (make(chan T, n))

  - Có buffer kích thước n.
  - send chỉ block khi buffer đầy, receive block khi buffer rỗng.
  - Dùng khi muốn tách nhịp producer/consumer (queue, worker pool).

  Ví dụ nhanh:

  ch1 := make(chan int)    // unbuffered
  ch2 := make(chan int, 2) // buffered

  ch2 <- 1 // không block
  ch2 <- 2 // không block
  // ch2 <- 3 // block nếu chưa có ai receive

  Câu chốt khi phỏng vấn:
  unbuffered = synchronization, buffered = decoupling + throughput.
