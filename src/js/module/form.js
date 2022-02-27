'uss strict';

// DOM TREEの構築が完了したら処理開始
document.addEventListener('DOMContentLoaded', () => {
  // DOM APIを利用してHTML要素を所得する
  const inputs = document.getElementsByTagName('input');
  const form = document.forms.namedItem('article-form');
  const saveBtn = document.querySelector('.article-form__save');
  const cancelBtn = document.querySelector('.article-form__cancel');
  const previewOpenBtn = document.querySelector('.article-form__open-preview');
  const previewCloseBtn = document.querySelector('.article-form__close-preview');
  const articleFormBody = document.querySelector('.article-form__body');
  const articleFormPreview = document.querySelector('.article-form__preview');
  const articleFormBodyTextArea = document.querySelector('.article-form__input--body');
  const articleFormPreviewTextArea = document.querySelector('.article-form__preview-body-contents');
  const csrfToken = document.getElementsByName('csrf')[0].content;
  const errors = document.querySelector('.article-form__errors');
  const errorTmpl = document.querySelector('.article-form__error-tmpl').firstElementChild;

  // 新規作成画面か編集画面かURLで判定
  const mode = { method: '', url: '' };
  if (window.location.pathname.endsWith('new')) {
    // newの場合HTTPメソッドはPOSTを利用
    mode.method = 'POST';
    // 作成リクエスト、及び戻るボタンの遷移先は「/」
    mode.url = '/articles';
  } else if (window.location.pathname.endsWith('edit')) {
    mode.method = 'PATCH';
    // 更新リクエスト、及び戻るボタンの遷移先は「/articles/:articleID」
    mode.url = `/articles/${window.location.pathname.split('/')[2]}`;
  }
  const { method, url } = mode;

  // input要素にフォーカスしてEnter押すとForm送信されるので、今回はEnterでForm送信されないよう制御
  for (let elm of inputs) {
    elm.addEventListener('keydown', event => {
      if (event.keyCode && event.keyCode === 13) {
        // キーが押された際のデフォルト挙動をキャンセル
        event.preventDefault();
        return false;
      }
    });
  }

  // プレビューを開くイベントを設定
  previewOpenBtn.addEventListener('click', event => {
    // formの「本文」に入力された内容をプレビューにコピーする
    articleFormPreviewTextArea.innerHTML = md.render(articleFormBodyTextArea.value);
    //入力フォームを非表示に プレビューを表示
    articleFormBody.style.display = 'none';
    articleFormPreview.style.display = 'grid'
  });

  // プレビューを閉じるイベントを設定
  previewCloseBtn.addEventListener('click', event => {
    articleFormBody.style.display = 'grid';
    articleFormPreview.style.display = 'none'
  });

  //前のページに戻るイベントを設定
  cancelBtn.addEventListener('click', event => {
    // button 要素クリック時のデフォルトの挙動をキャンセル
    event.preventDefault();

    //URLを指定して画面を遷移させる
    window.location.href = url;
  });

  // 保存処理イベント設定
  saveBtn.addEventListener('click', event => {
    // button 要素クリック時のデフォルトの挙動をキャンセル
    event.preventDefault();
    //前回のエラーメッセージを削除
    errors.innerHTML = null;
    // フォームの内容を取得
    const fd = new FormData(form);
    let status;
    // fetch APIを利用してリクエストを送信
    fetch(`/api${url}`, {
      method: method,
      headers: { 'X-CSRF-Token': csrfToken },
      body: fd
    })
      .then(res => {
        status = res.status;
        return res.json();
      })
      .then(body => {
        console.log(JSON.stringify(body));
        if (status === 200) {
          // 成功時は一覧に繊維
          window.location.href = url;
        }
        if (body.ValidationErrors) {
          showErrors(body.ValidationErrors);
        }
      })
      .catch(err => console.error(err));
  });
  const showErrors = messages => {
    //引数の値が配列であることを確認します。
    if (Array.isArray(messages) && messages.length != 0) {
      //複数メッセージを格納するためのフラグメントを作成
      const fragment = document.createDocumentFragment();
      //メッセージのループ処理
      messages.forEach(messages => {
        //単一メッセージ格納のためのフラグメントを作成
        const frag = document.createDocumentFragment();
        // テンプレートをクローンしてフラグメントに追加する。
        frag.appendChild(errorTmpl.cloneNode(true));
        // エラー要素にメッセージをセットする
        frag.querySelector('.article-form__error').innerHTML = messages;
        // エラー要素を親フラグメントに追加
        fragment.appendChild(frag);
      });
      // エラーメッセージの表示エリア（要素）にメッセージを追加
      errors.appendChild(fragment);
    }
  };
});