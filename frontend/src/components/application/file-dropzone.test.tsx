import { fireEvent, render, screen } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { FileDropzone } from './file-dropzone'
it('accepts an actual dropped file and supports browsing and replacement', () => {
 const select = vi.fn(), error = vi.fn(), file = new File(['test'], 'students.xlsx')
 const { container, rerender } = render(<FileDropzone file={null} accept=".xlsx" label="Drop a roster" help="500 rows" onSelect={select} onError={error} />)
 const input = container.querySelector('input')!, click = vi.spyOn(input, 'click')
 fireEvent.click(screen.getByRole('button', { name: 'Browse files' })); expect(click).toHaveBeenCalledOnce()
 fireEvent.drop(container.querySelector('.file-dropzone')!, { dataTransfer: { files: [file] } }); expect(select).toHaveBeenCalledWith(file)
 fireEvent.drop(container.querySelector('.file-dropzone')!, { dataTransfer: { files: [file, file] } }); expect(error).toHaveBeenCalledWith('Choose one file at a time.')
 rerender(<FileDropzone file={file} accept=".xlsx" label="Drop a roster" help="500 rows" onSelect={select} onError={error} />)
 expect(screen.getByText(/students.xlsx/)).toBeInTheDocument();fireEvent.click(screen.getByRole('button', { name: 'Remove selected file' })); expect(select).toHaveBeenLastCalledWith(null)
})
it('ignores drops during upload', () => {
 const select = vi.fn(), error = vi.fn(); const { container } = render(<FileDropzone file={null} accept=".xlsx" label="Drop a roster" help="500 rows" busy onSelect={select} onError={error} />)
 fireEvent.drop(container.querySelector('.file-dropzone')!, { dataTransfer: { files: [new File(['x'], 'x.xlsx')] } }); expect(select).not.toHaveBeenCalled();expect(screen.getByRole('button', { name: 'Browse files' })).toBeDisabled();expect(screen.getByRole('status')).toHaveTextContent('Processing')
})

it('keeps an image preview inside the uploader and prevents pending selection changes', () => {
 const select = vi.fn(), file = new File(['test'], 'catalog.png', {type:'image/png'})
 const {container} = render(<FileDropzone file={file} accept="image/png" label="Drop a catalog image" help="512 KiB" busy onSelect={select} onError={vi.fn()} preview={<img alt="Selected catalog preview" />} />)
 expect(container.querySelector('.file-dropzone')).toContainElement(screen.getByRole('img',{name:'Selected catalog preview'}))
 expect(screen.getByRole('button',{name:'Replace file'})).toBeDisabled()
 expect(screen.getByRole('button',{name:'Remove selected file'})).toBeDisabled()
 fireEvent.drop(container.querySelector('.file-dropzone')!,{dataTransfer:{files:[file]}})
 expect(select).not.toHaveBeenCalled()
 expect(screen.getByRole('status')).toHaveTextContent('Processing file')
})
