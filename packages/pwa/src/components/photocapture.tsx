import { CameraAltOutlined } from "@mui/icons-material";
import { Box, Button, Stack, styled } from "@mui/material"
import { useTranslation } from "react-i18next";

export type Props = {
    image?: File
    setImage: (file: File) => void
}

const VisuallyHiddenInput = styled('input')`
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    height: 1px;
    overflow: hidden;
    position: absolute;
    bottom: 0;
    left: 0;
    white-space: nowrap;
    width: 1px;
`;

export default function PhotoCapture({image, setImage}: Props) {
    const { t } = useTranslation()

    return (
        <Stack spacing={2}>
            { image && (
                <Box sx={{ display: 'flex', width: '100%', direction: 'column'}}>
                    <img src={URL.createObjectURL(image)} width='100%' />
                </Box>
            )}
            <Button
                component='label'
                variant='outlined'
                fullWidth
                startIcon={<CameraAltOutlined />}
            >
                {t(image ? 'action:retakePhoto' : 'action:takePhoto')}
                <VisuallyHiddenInput
                    type='file'
                    accept='image/*'
                    capture
                    onChange={(e) => {
                        if (!e.target.files || e.target.files.length === 0) {
                            return
                        }
                        setImage(e.target.files[0])
                    }}
                />
            </Button>
        </Stack>
    )
}