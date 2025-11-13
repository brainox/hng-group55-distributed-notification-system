import asyncio
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base
from app.core.config import settings

Base = declarative_base()

# --- SQLAlchemy setup ---
engine = create_async_engine(settings.database_url, echo=False, future=True)
async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)


async def init_db():
    """Initialize the database connection and wait until ready."""
    await wait_for_db(engine)
    print("✅ Database initialized successfully")


async def wait_for_db(engine):
    """Wait until the database is ready (retry logic)."""
    for attempt in range(10):
        try:
            async with engine.begin() as conn:
                await conn.run_sync(lambda _: None)
            print("✅ Database connection successful")
            return
        except Exception:
            print(f"⏳ Waiting for database to be ready... (attempt {attempt+1}/10)")
            await asyncio.sleep(3)
    raise Exception("❌ Could not connect to the database after multiple attempts")
