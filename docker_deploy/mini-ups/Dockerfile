FROM python:3.9-slim

WORKDIR /app

COPY . /app/

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

ENV PYTHONUNBUFFERED 1

# setup database

CMD ["bash", "-c", "python manage.py migrate && python manage.py collectstatic --noinput && gunicorn myproject.wsgi:application --bind 0.0.0.0:8000"]