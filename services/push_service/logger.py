#!/usr/bin/env python3
import os
import logging
from datetime import datetime
from logging.handlers import RotatingFileHandler


def setup_logger() -> logging.Logger:

    os.makedirs("logs", exist_ok=True)

    class CustomFormatter(logging.Formatter):
        def format(self, record):
            record.name = record.name.ljust(0)
            return super().format(record)

    logger = logging.getLogger("push_service")
    logger.setLevel(logging.DEBUG)

    file_handler = RotatingFileHandler(
        filename=f"logs/pushservice_{datetime.now().day:02d}_{datetime.now().month:02d}_{datetime.now().year}.log",
        maxBytes=10_000_000,
        backupCount=5,
        encoding="utf-8",
    )
    file_handler.setLevel(logging.DEBUG)
    file_formatter = CustomFormatter(
        fmt="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    file_handler.setFormatter(file_formatter)

    logger.addHandler(file_handler)

    return logger


Logger = setup_logger()
if __name__ == "__main__":
    Logger.info("Logger initialized successfully.")
