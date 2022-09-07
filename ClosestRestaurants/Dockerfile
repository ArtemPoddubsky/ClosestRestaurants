FROM golang:1.18 AS builder
COPY . /ClosestRestaurants
WORKDIR /ClosestRestaurants
RUN make all

FROM golang:1.18
COPY --from=builder /ClosestRestaurants ./
CMD ./app