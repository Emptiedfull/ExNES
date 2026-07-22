import os 
import sys 

root = sys.argv[1] if len(sys.argv) > 1 else "."

for dirpath,ignore,items in os.walk(root,topdown=False):
   
    for name in items:
        lower = name.lower()
        if name == lower:
            continue
        src = os.path.join(dirpath, name)
        dst = os.path.join(dirpath, lower)
      
        tmp = os.path.join(dirpath, lower + ".tmp_rename")
        os.rename(src, tmp)
        os.rename(tmp, dst)