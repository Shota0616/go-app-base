// ローカルストレージからキーに対応する値を取得し、存在しなければ初期値を返す関数
const getLocalStorageValue = (key, initValue) => {
// 指定したキーでローカルストレージから値を取得
const item = localStorage.getItem(key);

// 値が存在すればその値を返し、存在しなければ初期値を返す
return item ? item : initValue;
};

// カスタムフックを定義。指定したキーと初期値に基づいて、ローカルストレージと状態を同期する
export const useLocalStorage = (key, initValue) => {
// useStateで状態を管理。初期値はローカルストレージの値か、存在しなければ指定された初期値
const [value, setValue] = useState(() =>
    getLocalStorageValue(key, initValue)
);

// useEffectを使って、ローカルストレージの変更を監視
useEffect(() => {
    // ローカルストレージが変更されたときに呼び出されるコールバック関数を定義
    const callback = (event) => {
    // 変更されたキーが現在監視しているキーと一致するかを確認
    if (event.key === key) {
        // 一致すれば、ローカルストレージの新しい値で状態を更新
        setValue((value) => localStorage.getItem(key) ?? value);
    }
    };

    // 'storage'イベントリスナーを追加。これは同一ドメイン内の異なるタブでストレージが変更されたときに発火
    window.addEventListener('storage', callback);

    // クリーンアップ関数を返す。コンポーネントがアンマウントされたときにイベントリスナーを削除
    return () => {
    window.removeEventListener('storage', callback);
    };
}, [key]); // keyが変更されるたびに再実行

// ローカルストレージと状態を同時に更新する関数を定義
const setLocalStorageValue = useCallback(
    // setStateActionが関数の場合と文字列の場合で処理を分ける
    (setStateAction) => {
    // setStateActionが関数の場合、現在の値を引数にして新しい値を計算
    const newValue =
        setStateAction instanceof Function
        ? setStateAction(value)
        : setStateAction;

    // ローカルストレージを新しい値で更新
    localStorage.setItem(key, newValue);
    // storageイベント発火
    window.dispatchEvent(new Event("storage"));
    // Reactの状態も新しい値で更新
    setValue(() => newValue);
    },
    [key, value] // keyとvalueが変わるときにこの関数を再生成
);

// 現在の値と更新関数を返す
return [value, setLocalStorageValue];
};

// ローカルストレージが変更されたときに呼び出されるコールバック関数
const callback = (event) => {
// 変更されたキーが指定されたキーと一致する場合
if (event.key === key) {
    // ローカルストレージの新しい値で状態を更新する
    setValue((value) => localStorage.getItem(key) ?? value);
}
};
