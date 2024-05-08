function previewImage() {
    const previewContainer = document.getElementById('previewContainer');
    const file = document.getElementById('imageInput').files[0];
    const reader = new FileReader();

    reader.onload = function(e) {
        // 清除旧的预览
        previewContainer.innerHTML = '';
        // 创建一个新的 img 元素并设置 src 属性
        const img = new Image();
        img.src = e.target.result;
        img.style.maxWidth = '100%'; // 确保图片不会超过容器宽度
        img.style.borderRadius = '5px'; // 图片圆角
        previewContainer.appendChild(img);
    };

    if (file) {
        reader.readAsDataURL(file);
    }
}

// 保留原有的 uploadImage 函数...


function uploadImage() {
    const imageInput = document.getElementById('imageInput');
    const loadingText = document.getElementById('loadingText');
    const resultText = document.getElementById('resultText');

    if (imageInput.files.length === 0) {
        alert('请选择一张图片！');
        return;
    }

    // 显示加载指示器
    loadingText.style.display = 'block';
    // 确保结果文本框隐藏
    resultText.style.display = 'none';

    const formData = new FormData();
    formData.append('image', imageInput.files[0]);

    fetch('http://127.0.0.1:8080/api/v1/upload/image', {
            method: 'POST',
            body: formData,
        })
        .then(response => response.json())
        .then(data => {
            // 隐藏加载指示器
            loadingText.style.display = 'none';
            // 显示结果文本框并设置文本
            resultText.style.display = 'block';
            resultText.innerText = data.data;
        })
        .catch(error => {
            console.error('Error:', error);
            // 出错时也需要更新界面
            loadingText.style.display = 'none';
            resultText.style.display = 'block';
            resultText.innerText = '发生错误，请重试';
        });
}