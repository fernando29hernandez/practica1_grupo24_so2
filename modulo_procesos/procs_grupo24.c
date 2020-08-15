#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/uaccess.h>
#include <linux/slab.h>
#include <linux/fs.h>
#include <linux/sysinfo.h>
#include <linux/seq_file.h>
#include <linux/slab.h>
#include <linux/mm.h>
#include <linux/swap.h>

static int open_action(struct seq_file *m, void *v)
{

    struct task_struct *task;

    seq_printf(m, "Pid,Nombre_Del_Proceso,Estado,Usuario,Padre Pid,Memoria\n");

    for_each_process(task)
    {

        if (task->mm)
        {
            if (task->parent)
            {

                seq_printf(m, "%d,%s,%d,%u,%d,%d\n", task->pid, task->comm, task->state, task->cred->uid, task->parent->pid, task->mm->total_vm);
            }
            else
            {
                seq_printf(m, "%d,%s,%d,%u,%s,%d\n", task->pid, task->comm, task->state, task->cred->uid, "NULL", task->mm->total_vm);
            }
        }
        else
        {
            if (task->parent)
            {
                seq_printf(m, "%d,%s,%d,%u,%d,%d\n", task->pid, task->comm, task->state, task->cred->uid, task->parent->pid, 0);
            }
            else
            {
                seq_printf(m, "%d,%s,%d,%u,%s,%d\n", task->pid, task->comm, task->state, task->cred->uid, "NULL", 0);
            }
        }
    }

    return 0;
}
ssize_t write_proc(struct file *filp, const char *buf, size_t count, loff_t *offp)
{
    return 0;
}

int open_proc(struct inode *inode, struct file *file)
{
    return single_open(file, open_action, NULL);
}

static struct file_operations proc_fops = {
    read : seq_read,
    write : write_proc,
    open : open_proc,
    release : single_release,
    llseek : seq_lseek
};

static int __init ejemplo_init(void)
{
    printk(KERN_INFO "Hola mundo, somos el grupo 24 y este es el monitor de procesos \n");

    struct proc_dir_entry *entry;
    entry = proc_create("procs_grupo24", 0777, NULL, &proc_fops);

    if (!entry)
    {
        return -1;
    }

    return 0;
}

static void __exit ejemplo_exit(void)
{
    printk(KERN_INFO "Sayonara mundo, somos el grupo 24  y este fue el monitor de procesos\n");
    remove_proc_entry("procs_grupo24", NULL);
}

module_init(ejemplo_init);
module_exit(ejemplo_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Fernando Hernandez");
MODULE_DESCRIPTION("CPU");
MODULE_VERSION("0.01");

