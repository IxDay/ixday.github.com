FROM debian

WORKDIR /mnt

COPY requirements.txt /tmp

RUN apt-get update
RUN apt-get install -y python-dev python-pip git

RUN pip install -r /tmp/requirements.txt

RUN git clone https://github.com/IxDay/pelican-chunk /tmp/pelican-chunk
RUN pelican-themes -i /tmp/pelican-chunk

ENTRYPOINT ["fab"]
CMD ["serve"]
EXPOSE 8000
