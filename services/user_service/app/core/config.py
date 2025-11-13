from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    app_name: str = "user_service"
    database_url: str
    jwt_secret: str
    jwt_algorithm: str = "HS256"
    access_token_expire_minutes: int = 30

    class Config:
        env_file = ".env"

settings = Settings()
