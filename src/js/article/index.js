'uss strict';

// DOM TREEの構築が完了したら処理開始
document.addEventListener('DOMContentLoaded', () => {
  // DOM APIを利用してHTML要素を所得する
  const deleteBtns = document.querySelectorAll('.articles__item-delete')
  const csrfToken = document.getElementsByName('csrf')[0].content;

  const deleteArticle = id => {
    let statusCode;
    fetch(`/${id}`, {
      method: 'DELETE',
      headers: { 'X-CSRF-Token': csrfToken }
    })
      .then(res => {
        statusCode = res.status;
        return res.json();
      })
      .then(data => {
        console.log(JSON.stringify(data));
        if (statusCode === 200) {
          //削除に成功したら画面から記事のHTMLを削除
          document.querySelector(`.articles__item-${id}`).remove();
        }
      })
      .catch(err => console.error(err));
  };

  //削除ボタンそれぞれに対し、イベントリスナーを設定
  for (let elm of deleteBtns) {
    elm.addEventListener('click', event => {
      event.preventDefault();
      deleteArticle(elm.dataset.id);
    });
  }
});