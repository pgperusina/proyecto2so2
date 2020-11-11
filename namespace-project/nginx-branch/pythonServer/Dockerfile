FROM debian
LABEL maintainer="PGGP"
RUN apt-get -y update
RUN apt-get install -y python3-venv python3-pip

WORKDIR /code
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt
COPY . .
CMD ["python3", "-u", "pythonApp.py"]
# ENTRYPOINT ["python3", "pythonApp.py"]