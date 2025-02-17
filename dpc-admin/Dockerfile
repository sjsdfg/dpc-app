FROM ruby:2.7.2-alpine AS ruby_builder

# Install build dependencies
RUN apk add --no-cache postgresql-dev && \
    apk add --no-cache libsodium-dev && \
    apk add --no-cache shared-mime-info && \
    apk add --no-cache --virtual build-deps alpine-sdk npm tzdata

# Set the working directory
RUN mkdir /dpc-admin
WORKDIR /dpc-admin

# Copy over the files needed to fetch dependencies
COPY /dpc-admin/Gemfile /dpc-admin/Gemfile.lock /dpc-admin/
COPY /dpc-admin/package.json /dpc-admin/package-lock.json /dpc-admin/

# Install the website dependencies
RUN gem install bundler --no-document && bundle install && npm install

# Run bundler audit
RUN bundle exec bundle audit update && bundle exec bundle audit check --ignore CVE-2015-9284

# Copy the code, test the app, and build the assets pipeline
COPY /dpc-admin /dpc-admin

# TODO: Determine if we have to remove /bin/yarn
RUN rm -rf /bin/yarn
RUN RAILS_ENV=production WEBPACKER_NODE_MODULES_BIN_PATH="node_modules/.bin" SECRET_KEY_BASE=dummy bundle exec rake assets:precompile --trace

# Clean up from the build
RUN rm -rf /usr/local/bundle/cache/*.gem && \
    find /usr/local/bundle/gems/ -name "*.c" -delete && \
    find /usr/local/bundle/gems/ -name "*.o" -delete


FROM ruby:2.7.2-alpine

# Install app dependencies
RUN apk add --no-cache postgresql-client nodejs tzdata libsodium-dev

# Set the working directory
RUN mkdir /dpc-admin
WORKDIR /dpc-admin

# Copy the compiled Ruby app
COPY --from=ruby_builder /usr/local/bundle/ /usr/local/bundle/
COPY --from=ruby_builder /dpc-admin /dpc-admin

# Declare the entrypoint shell script
ENTRYPOINT ["./docker/entrypoint.sh"]

# Default to running the rails server
CMD ["admin"]
