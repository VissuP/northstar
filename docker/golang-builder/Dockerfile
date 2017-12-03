FROM buildpack-deps:jessie-scm

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		patch \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	&& rm -rf /var/lib/apt/lists/*

# Download bootstrap binary
ENV GOLANG_VERSION 1.9
ENV GOLANG_DOWNLOAD_URL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 d70eadefce8e160638a9a6db97f7192d8463069ab33138893ad3bf31b0650a79

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256 golang.tar.gz" | sha256sum -c - \
	&& tar -C $HOME -xzf golang.tar.gz \
	&& rm golang.tar.gz

# Download source
ENV GOLANG_DOWNLOAD_URL https://storage.googleapis.com/golang/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 a4ab229028ed167ba1986825751463605264e44868362ca8e7accc8be057e993

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& mv $HOME/go $HOME/go1.4 \
	&& rm golang.tar.gz

# Apply cover patches
WORKDIR /usr/local/go/src
COPY patches/testing_cover_go.patch /usr/local/go/src/testing_cover_go.patch
RUN patch -p1 < testing_cover_go.patch
RUN rm testing_cover_go.patch
COPY patches/cmd_cover_cover_go.patch /usr/local/go/src/cmd_cover_cover_go.patch
RUN patch -p1 < cmd_cover_cover_go.patch
RUN rm cmd_cover_cover_go.patch
COPY patches/cmd_cover_func_go.patch /usr/local/go/src/cmd_cover_func_go.patch
RUN patch -p1 < cmd_cover_func_go.patch
RUN rm cmd_cover_func_go.patch
COPY patches/cover2clover.go /usr/local/go/src/cmd/cover/cover2clover.go

# Build from source
# See https://dave.cheney.net/2013/06/04/how-go-uses-go-to-build-itself
RUN ./all.bash || true; \
	rm -rf $HOME/go1.4

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY go-wrapper /usr/local/bin/
