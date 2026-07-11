document.addEventListener('DOMContentLoaded', () => {
    const shareBtns = document.querySelectorAll('.share-btn');
    shareBtns.forEach(btn => {
        btn.addEventListener('click', async () => {
            const id = btn.getAttribute('data-share-id');
            const name = btn.getAttribute('data-share-name');
            const url = window.location.origin + '/files/' + id + '?download=true';

            const shareData = {
                title: 'PLUS Drop',
                text: name + ' 파일을 공유합니다.',
                url: url
            };

            if (navigator.share) {
                try {
                    await navigator.share(shareData);
                } catch (err) {
                    if (err.name !== 'AbortError') {
                        console.error('Share failed:', err);
                    }
                }
            } else {
                try {
                    await navigator.clipboard.writeText(url);
                    alert('공유 링크가 클립보드에 복사되었습니다.');
                } catch (err) {
                    alert('링크 복사에 실패했습니다. 수동으로 복사해주세요: ' + url);
                }
            }
        });
    });
});
