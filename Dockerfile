FROM scratch

ARG VERSION
ARG ARCH
ARG OS=linux

COPY ./build/bin/${OS}_${ARCH}_${VERSION}/pc-proxy /

# Command to run
ENTRYPOINT ["/pc-proxy"]
