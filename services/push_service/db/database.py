#!/usr/bin/env python3
import psycopg2
from decouple import config
from sqlmodel import SQLModel
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker


DATABASE_URL = "sqlite:///./push.db"

engine = create_engine(
    DATABASE_URL, connect_args={"check_same_thread": False}, echo=False
)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


def init_db():
    SQLModel.metadata.create_all(bind=engine)  # type: ignore
