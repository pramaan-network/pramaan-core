FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY pramaand /usr/bin/pramaand

RUN chmod +x /usr/bin/pramaand

RUN mkdir -p /root/.pramaand

EXPOSE 26656 26657

CMD ["pramaand", "start"]
