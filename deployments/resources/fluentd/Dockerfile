FROM fluent/fluentd:edge-debian

USER root

RUN apt-get update && apt-get install -y netcat-openbsd
RUN gem install fluent-plugin-grafana-loki

USER fluent
