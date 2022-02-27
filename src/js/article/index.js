'uss strict';

// DOM TREEの構築が完了したら処理開始
document.addEventListener('DOMContentLoaded', () => {
  // DOM APIを利用してHTML要素を所得する
  const deleteBtns = document.querySelectorAll('.articles__item-delete')
  const moreBtn = document.querySelector('.page__more');
  const articles = document.querySelector('.articles');
  const articleTmpl = document.querySelector('.articles__item-tmpl').firstElementChild;

  const csrfToken = document.getElementsByName('csrf')[0].content;

  const deleteArticle = id => {
    let statusCode;
    fetch(`/api/articles/${id}`, {
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

  moreBtn.addEventListener('click', event => {
    event.preventDefault();
    //もっと見るボタンのカスタムデータ属性からカーソルの値を取得
    const cursor = moreBtn.dataset.cursor;

    //cursorが取得できない場合や0以下の数値だった場合はもっと見るボタンを画面から削除して処理を終了する
    if (!cursor || cursor <= 0) {
      moreBtn.remove();
      return;
    }

    //FETCHAPIを使用して非同期リクエストを実行する
    let statusCode;
    fetch(`/api/articles?cursor=${cursor}`)
      .then(res => {
        statusCode = res.status;
        return res.json();
      })
      .then(data => {
        console.log(JSON.stringify(data));
        //リクエストに成功し記事一覧データが配列で帰ってきた場合
        if (statusCode == 200 && Array.isArray(data)) {
          //表示する記事がこれ以上存在しない場合はもっと見るボタンを画面から削除して処理終了
          if (data.length == 0) {
            moreBtn.remove();
            return;
          }
          //記事のHTML要素をまとめるためのフラグメントを作成する（記事のリスト）
          const fragment = document.createDocumentFragment();
          //記事一覧データをループ処理
          data.forEach(article => {
            //個々の記事のHTML要素を格納するフラグメントを作成する（個別記事）
            const frag = document.createDocumentFragment();
            //記事のHTML要素のテンプレートからクローンを作成しフラグメントの子要素として追加
            frag.appendChild(articleTmpl.cloneNode(true));
            //記事の各HTML要素に対しクラス・属性値、テキストを設定する
            frag.querySelector('article').classList.add(`articles__item-${article.id}`);
            frag.querySelector('.articles__item').setAttribute('href', `/articles/${article.id}`);
            frag.querySelector('.articles__item-title').textContent = article.title;
            frag.querySelector('.articles__item-date').textContent = article.created.split('T')[0]; //+年-月-日のみを抽出
            //DeleteButtonに対してカスタムデータ属性やイベントリスナーを設定する
            const deleteBtnElm = frag.querySelector('.articles__item-delete');
            deleteBtnElm.dataset.id = article.id;
            deleteBtnElm.addEventListener('click', event => {
              event.preventDefault();
              deleteArticle(article.id);
            });
            //記事リストのフラグメントの子要素として個別記事のフラグメントを追加する
            fragment.appendChild(frag);
          });
          //もっと見るボタンのカスタムデータ属性に設定してあるカーソルの値を更新する
          moreBtn.dataset.cursor = data[data.length - 1].id;
          //記事一覧のHTML要素の子要素に記事リストのフラグメントを追加して画面に表示する
          articles.appendChild(fragment);
        }
      })
      .catch(err => console.error(err));
  });
});