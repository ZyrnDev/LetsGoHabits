FROM habits_backend:latest

ARG EXEC_APP_NAME
ENV EXEC_APP_NAME=${EXEC_APP_NAME:-development}

CMD $EXEC_APP_NAME