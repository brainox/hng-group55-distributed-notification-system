from logging.config import fileConfig
from sqlalchemy import pool
from sqlalchemy.ext.asyncio import create_async_engine
from alembic import context
import asyncio
import os

# Import your Base metadata
from app.core.db import Base
from app.models.user import User
from app.models.preference import Preference
from app.models.push_token import PushToken

# this is the Alembic Config object
config = context.config
fileConfig(config.config_file_name)
target_metadata = Base.metadata

# Get DATABASE_URL from environment variable, fallback to alembic.ini
DATABASE_URL = os.getenv("DATABASE_URL") or config.get_main_option("sqlalchemy.url")

def run_migrations_offline():
    context.configure(
        url=DATABASE_URL,
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )
    with context.begin_transaction():
        context.run_migrations()

def do_run_migrations(connection):
    context.configure(connection=connection, target_metadata=target_metadata)
    with context.begin_transaction():
        context.run_migrations()

async def run_migrations_online():
    connectable = create_async_engine(DATABASE_URL, poolclass=pool.NullPool)

    async with connectable.connect() as connection:
        await connection.run_sync(do_run_migrations)

if context.is_offline_mode():
    run_migrations_offline()
else:
    asyncio.run(run_migrations_online())
