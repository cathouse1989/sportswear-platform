from pathlib import Path
p = Path(r'D:\code\sportswear\sportswear-platform\internal\services\cms_service.go')
text = p.read_text(encoding='utf-8')
start = text.find('// HeroSlide')
end = text.find('// DeletePage')
print('start', start, 'end', end)
if start>=0:
    print(text[start:start+120])
