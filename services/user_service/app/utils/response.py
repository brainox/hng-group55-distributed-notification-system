from typing import Any


def response(success: bool, data:Any=None, message="", error=None, meta=None):
    if meta is None:
        meta = {
            "total": 0, "limit": 0, "page": 0, "total_pages": 0,
            "has_next": False, "has_previous": False
        }
    return {
        "success": success,
        "data": data,
        "error": error,
        "message": message,
        "meta": meta
    }
