FROM discoenv/javabase

USER root
VOLUME ["/etc/iplant/de"]

COPY conf/main/logback.xml /
COPY target/publish-standalone.jar /

ARG git_commit=unknown
ARG buildenv_git_commit=unknown
ARG version=unknown
LABEL org.iplantc.de.publish.git-ref="$git_commit" \
      org.iplantc.de.publish.version="$version" \
      org.iplantc.de.buildenv.git-ref="$buildenv_git_commit"

RUN ln -s "/opt/jdk/bin/java" "/bin/publish"
ENTRYPOINT ["publish", "-Dlogback.configurationFile=/etc/iplant/de/logging/publish-logging.xml", "-cp", ".:publish-standalone.jar:/", "publish.core"]
CMD ["--help"]
