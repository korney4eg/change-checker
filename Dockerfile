FROM ruby:3-buster

RUN apt-get update && apt-get install -y \
  gcc\ 
  bash\ 
  cmake\ 
  git \
  && rm -rf /var/lib/apt/lists/*

# install both bundler 1.x and 2.x
RUN gem install bundler -v "~>1.0" && gem install bundler jekyll

# RUN apk --no-cache add ca-certificates
ADD https://github.com/korney4eg/change-checker/releases/download/v0.2.1/change-checker /change-checker
RUN chmod +x /change-checker
CMD ["/change-checker"]