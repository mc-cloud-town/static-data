import json
import os
from pathlib import Path
from typing import Dict, List, Optional, TypedDict

import requests

BASE_API_URL = os.getenv("BASE_API_URL", "https://example.com/api")
MEMBERS_API_URL = f"{BASE_API_URL}/members"


class Member(TypedDict):
    uuid: str
    name: str
    introduction: Optional[str]


# { [group: str]: Member[] }
def get_members() -> Dict[str, List[Member]]:
    return requests.get(MEMBERS_API_URL).json()


OUTPUT_DIR = Path("output")
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
(OUTPUT_DIR / "member.json").write_text(json.dumps(get_members(), separators=(",", ":")))
