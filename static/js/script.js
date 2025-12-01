// Konfirmasi sebelum menghapus
document.addEventListener('DOMContentLoaded', function() {
    const deleteButtons = document.querySelectorAll('.btn-danger');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            if (!confirm('Apakah Anda yakin ingin menghapus produk ini?')) {
                e.preventDefault();
            }
        });
    });
});