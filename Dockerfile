FROM google/debian:wheezy

MAINTAINER Stefan Arentz <stefan@arentz.ca>

ADD moz-mockmyid-api /usr/local/bin/moz-mockmyid-api
ADD moz-mockmyid-api.sh /usr/local/bin/moz-mockmyid-api.sh

EXPOSE 8080

CMD ["/usr/local/bin/moz-mockmyid-api.sh"]
