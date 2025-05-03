# Use a full-featured base with bash, curl, etc.
FROM ubuntu:24.04

# install common tools + HTTPS support
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  bash \
  ca-certificates \
  curl \
  less \
  vim \
  && rm -rf /var/lib/apt/lists/*

# copy your CLI into place
COPY ti /usr/local/bin/ti
RUN chmod +x /usr/local/bin/ti

# ensure it's on the PATH
ENV PATH="/usr/local/bin:${PATH}"

# default to an interactive bash shell
# if you pass "ti ..." as args, it'll run your CLI instead
CMD ["bash"]
